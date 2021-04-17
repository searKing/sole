// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/searKing/sole/internal/pkg/version"
	"github.com/searKing/sole/pkg/consul"
	vipergetter "github.com/searKing/sole/pkg/viper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	viperhelper "github.com/searKing/golang/third_party/github.com/spf13/viper"

	viper_ "github.com/searKing/sole/api/protobuf-spec/v1/viper"
	"github.com/searKing/sole/pkg/crypto/pasta"
	redis_ "github.com/searKing/sole/pkg/database/redis"
	"github.com/searKing/sole/pkg/database/sql"
	"github.com/searKing/sole/pkg/logs"
	"github.com/searKing/sole/pkg/opentrace"
)

//go:generate go-option -type=Config
type Config struct {
	GetViper func() *viper.Viper // If set, overrides params below
	proto    *viper_.ViperProto

	Logs       *logs.Config
	OpenTracer *opentrace.Config

	KeyCipher *pasta.Config
	Sql       *sql.Config
	Redis     *redis_.Config
	Consul    *consul.Config
}

type completedConfig struct {
	*Config

	// for Complete Only
	completeError error
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// NewConfig returns a Config struct with the default values
func NewConfig() *Config {
	return &Config{
		proto:      NewDefaultViperProto(),
		Logs:       logs.NewViperConfig(vipergetter.GetViper("log", version.ServiceName)),
		OpenTracer: opentrace.NewViperConfig(vipergetter.GetViper("tracing", version.ServiceName)),
		KeyCipher:  pasta.NewViperConfig(vipergetter.GetViper("secret", version.ServiceName)),
		Sql:        sql.NewViperConfig(vipergetter.GetViper("database", version.ServiceName)),
		Redis:      redis_.NewViperConfig(vipergetter.GetViper("redis", version.ServiceName)),
		Consul:     consul.NewViperConfig(vipergetter.GetViper("consul", version.ServiceName)),
	}
}

// NewViperConfig returns a Config struct with the global viper instance
// key representing a sub tree of this instance.
// NewViperConfig is case-insensitive for a key.
func NewViperConfig(getViper func() *viper.Viper) *Config {
	c := NewConfig()
	c.GetViper = getViper
	return c
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to ApplyOptions, do that first. It's mutating the receiver.
// ApplyOptions is called inside.
func (c *Config) Complete() CompletedConfig {
	if err := c.loadViper(); err != nil {
		return CompletedConfig{&completedConfig{
			Config:        c,
			completeError: err,
		}}
	}
	return CompletedConfig{&completedConfig{Config: c}}
}

// New creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
// New usually called after Complete
func (c completedConfig) New(ctx context.Context) (*Provider, error) {
	var sqlDB *sqlx.DB

	if err := c.Logs.Complete().Apply(); err != nil {
		return nil, err
	}
	if closer, err := c.OpenTracer.Complete().Apply(); err != nil {
		go func() {
			select {
			case <-ctx.Done():
				err := closer.Close()
				if err != nil {
					logrus.WithError(err).Error("openTracing closed")
					return
				}
				logrus.Info("openTracing closed")
			}
		}()
		return nil, err
	}

	if c.Sql != nil {
		sqlDB = c.Sql.Complete().New(ctx)
	}

	p := &Provider{
		Proto:     c.proto,
		SqlDB:     sqlDB,
		KeyCipher: c.KeyCipher.Complete().New(),
		ctx:       ctx,
	}
	if c.Redis != nil {
		redis, err := c.Redis.Complete().New()
		if err != nil {
			return nil, err
		}
		p.Redis = redis
	}

	if c.Consul != nil {
		builder := c.Consul.Complete()
		register, err := builder.NewServiceRegister()
		if err != nil {
			return nil, err
		}
		p.ServiceRegister = register
		resolver, err := builder.NewServiceResolver()
		if err != nil {
			return nil, err
		}
		p.ServiceResolver = resolver
	}

	providerReloads.WithLabelValues(p.Proto.String()).Inc()
	go p.ReloadForever()
	return p, nil
}

// Apply set options and something else as global init, act likes New but without Config's instance
// Apply usually called after Complete
func (c completedConfig) Apply(ctx context.Context) error {
	provider, err := c.New(ctx)
	if err != nil {
		return err
	}
	InitGlobalProvider(provider)
	return nil
}

func (c *Config) loadViper() error {
	var v *viper.Viper
	if c.GetViper != nil {
		v = c.GetViper()
	}

	if c.proto == nil {
		c.proto = &viper_.ViperProto{}
	}
	if err := viperhelper.UnmarshalProtoMessageByJsonpb(v, c.proto); err != nil {
		logrus.WithError(err).Errorf("load logs config from viper")
		return err
	}
	return nil
}
