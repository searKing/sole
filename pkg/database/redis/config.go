// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
	"github.com/searKing/sole/pkg/protobuf"
)

type Config struct {
	KeyInViper string
	Viper      *viper.Viper // If set, overrides params below
	Redis
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
	return &Config{}
}

// NewViperConfig returns a Config struct with the global viper instance
// key representing a sub tree of this instance.
// NewViperConfig is case-insensitive for a key.
func NewViperConfig(key string) *Config {
	c := NewConfig()
	c.Viper = viper.GetViper()
	c.KeyInViper = key
	return c
}

// Validate checks Config and return a slice of found errs.
func (c *Config) Validate() []error {
	return nil
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
func (c completedConfig) New() (redis.UniversalClient, error) {
	if c.completeError != nil {
		return nil, c.completeError
	}
	redis_ := c.Redis
	var UniversalOptions redis.UniversalOptions
	UniversalOptions.Addrs = redis_.GetAddrs()
	UniversalOptions.DB = int(redis_.GetDb())
	UniversalOptions.Username = redis_.GetUsername()
	UniversalOptions.Password = redis_.GetPassword()
	UniversalOptions.SentinelPassword = redis_.GetSentinelPassword()
	UniversalOptions.MaxRetries = int(redis_.GetMaxRetries())
	UniversalOptions.MinRetryBackoff = protobuf.DurationOrDefault(redis_.GetMinRetryBackoff(), UniversalOptions.MinRetryBackoff, "min_retry_backoff")
	UniversalOptions.MaxRetryBackoff = protobuf.DurationOrDefault(redis_.GetMaxRetryBackoff(), UniversalOptions.MaxRetryBackoff, "max_retry_backoff")
	UniversalOptions.DialTimeout = protobuf.DurationOrDefault(redis_.GetDialTimeout(), UniversalOptions.DialTimeout, "dial_timeout")
	UniversalOptions.ReadTimeout = protobuf.DurationOrDefault(redis_.GetReadTimeout(), UniversalOptions.ReadTimeout, "read_timeout")
	UniversalOptions.WriteTimeout = protobuf.DurationOrDefault(redis_.GetWriteTimeout(), UniversalOptions.WriteTimeout, "write_timeout")
	UniversalOptions.PoolSize = int(redis_.GetPoolSize())
	UniversalOptions.MinIdleConns = int(redis_.GetMinIdleConns())
	UniversalOptions.MaxConnAge = protobuf.DurationOrDefault(redis_.GetMaxConnAge(), UniversalOptions.MaxConnAge, "max_conn_age")
	UniversalOptions.PoolTimeout = protobuf.DurationOrDefault(redis_.GetPoolTimeout(), UniversalOptions.PoolTimeout, "pool_timeout")
	UniversalOptions.IdleTimeout = protobuf.DurationOrDefault(redis_.GetIdleTimeout(), UniversalOptions.IdleTimeout, "idle_timeout")
	UniversalOptions.IdleCheckFrequency = protobuf.DurationOrDefault(redis_.GetIdleCheckFrequency(), UniversalOptions.IdleCheckFrequency, "idle_check_frequency")
	UniversalOptions.MaxRedirects = int(redis_.GetMaxRedirects())
	UniversalOptions.ReadOnly = redis_.GetReadOnly()
	UniversalOptions.RouteByLatency = redis_.GetRouteByLatency()
	UniversalOptions.RouteRandomly = redis_.GetRouteRandomly()
	UniversalOptions.MasterName = redis_.GetMasterName()
	return redis.NewUniversalClient(&UniversalOptions), nil
}

func (c *Config) loadViper() error {
	v := c.Viper
	if v != nil && c.KeyInViper != "" {
		v = v.Sub(c.KeyInViper)
	}

	if err := viper_.UnmarshalProtoMessageByJsonpb(v, &c.Redis); err != nil {
		logrus.WithError(err).Errorf("load redis config from viper")
		return err
	}
	return nil
}
