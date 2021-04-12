// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/jwalterweatherman"

	viperhelper "github.com/searKing/golang/third_party/github.com/spf13/viper"

	viper_ "github.com/searKing/sole/api/protobuf-spec/v1/viper"
	"github.com/searKing/sole/internal/pkg/version"
	"github.com/searKing/sole/pkg/crypto/pasta"
	redis_ "github.com/searKing/sole/pkg/database/redis"
	"github.com/searKing/sole/pkg/database/sql"
	"github.com/searKing/sole/pkg/logs"
	"github.com/searKing/sole/pkg/opentrace"
)

//go:generate go-option -type=Config
type Config struct {
	ConfigFile string

	proto *viper_.ViperProto

	Logs       *logs.Config
	OpenTracer *opentrace.Config

	KeyCipher *pasta.Config
	Sql       *sql.Config
	Redis     *redis_.Config
}

// NewConfig returns a Config struct with the default values
func NewConfig() *Config {
	return &Config{
		Logs:       logs.NewViperConfig("log"),
		OpenTracer: opentrace.NewViperConfig("tracing"),
		KeyCipher:  pasta.NewViperConfig("secret"),
		Sql:        sql.NewViperConfig("database"),
		Redis:      redis_.NewViperConfig("redis"),
	}
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to ApplyOptions, do that first. It's mutating the receiver.
// ApplyOptions is called inside.
func (o *Config) Complete(options ...ConfigOption) CompletedConfig {
	o.ApplyOptions(options...)
	o.installViperProtoOrDie()
	return CompletedConfig{&completedConfig{o}}
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
		proto:     c.proto,
		sqlDB:     sqlDB,
		keyCipher: c.KeyCipher.Complete().New(),
		ctx:       ctx,
	}
	if c.Redis != nil {
		redis, err := c.Redis.Complete().New()
		if err != nil {
			return nil, err
		}
		p.redis = redis
	}
	providerReloads.WithLabelValues(p.proto.String()).Inc()
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

// installViperProtoOrDie allows you to load config from default, config path、env and so on, but dies on failure.
func (c *Config) installViperProtoOrDie() {
	var v viper_.ViperProto
	jwalterweatherman.SetLogOutput(logrus.StandardLogger().Writer())
	jwalterweatherman.SetLogThreshold(jwalterweatherman.LevelWarn)

	if err := viperhelper.LoadGlobalConfig(&v, c.ConfigFile, version.ServiceName, NewDefaultViperProto()); err != nil {
		logrus.WithError(err).WithField("config_path", c.ConfigFile).Fatalf("load config")
	}
	c.proto = &v
}
