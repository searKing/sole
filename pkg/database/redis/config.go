// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"

	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
	"github.com/searKing/sole/pkg/protobuf"
)

type Config struct {
	Proto     Redis
	Validator *validator.Validate

	viper     *viper.Viper
	viperKeys []string
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
func NewViperConfig(v *viper.Viper, keys ...string) *Config {
	c := NewConfig()
	c.viper = v
	c.viperKeys = keys
	return c
}

// Validate checks Config.
func (c *completedConfig) Validate() error {
	return c.Validator.Struct(c)
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to ApplyOptions, do that first. It's mutating the receiver.
// ApplyOptions is called inside.
func (c *Config) Complete() CompletedConfig {
	if c.viper != nil {
		err := viper_.UnmarshalKeysViper(c.viper, c.viperKeys, &c.Proto)
		if err != nil {
			return CompletedConfig{&completedConfig{completeError: err}}
		}
	}
	if c.Validator == nil {
		c.Validator = validator.New()
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
	err := c.Validate()
	if err != nil {
		return nil, err
	}
	var UniversalOptions redis.UniversalOptions
	UniversalOptions.Addrs = c.Proto.GetAddrs()
	UniversalOptions.DB = int(c.Proto.GetDb())
	UniversalOptions.Username = c.Proto.GetUsername()
	UniversalOptions.Password = c.Proto.GetPassword()
	UniversalOptions.SentinelPassword = c.Proto.GetSentinelPassword()
	UniversalOptions.MaxRetries = int(c.Proto.GetMaxRetries())
	UniversalOptions.MinRetryBackoff = protobuf.DurationOrDefault(c.Proto.GetMinRetryBackoff(), UniversalOptions.MinRetryBackoff, "min_retry_backoff")
	UniversalOptions.MaxRetryBackoff = protobuf.DurationOrDefault(c.Proto.GetMaxRetryBackoff(), UniversalOptions.MaxRetryBackoff, "max_retry_backoff")
	UniversalOptions.DialTimeout = protobuf.DurationOrDefault(c.Proto.GetDialTimeout(), UniversalOptions.DialTimeout, "dial_timeout")
	UniversalOptions.ReadTimeout = protobuf.DurationOrDefault(c.Proto.GetReadTimeout(), UniversalOptions.ReadTimeout, "read_timeout")
	UniversalOptions.WriteTimeout = protobuf.DurationOrDefault(c.Proto.GetWriteTimeout(), UniversalOptions.WriteTimeout, "write_timeout")
	UniversalOptions.PoolSize = int(c.Proto.GetPoolSize())
	UniversalOptions.MinIdleConns = int(c.Proto.GetMinIdleConns())
	UniversalOptions.MaxConnAge = protobuf.DurationOrDefault(c.Proto.GetMaxConnAge(), UniversalOptions.MaxConnAge, "max_conn_age")
	UniversalOptions.PoolTimeout = protobuf.DurationOrDefault(c.Proto.GetPoolTimeout(), UniversalOptions.PoolTimeout, "pool_timeout")
	UniversalOptions.IdleTimeout = protobuf.DurationOrDefault(c.Proto.GetIdleTimeout(), UniversalOptions.IdleTimeout, "idle_timeout")
	UniversalOptions.IdleCheckFrequency = protobuf.DurationOrDefault(c.Proto.GetIdleCheckFrequency(), UniversalOptions.IdleCheckFrequency, "idle_check_frequency")
	UniversalOptions.MaxRedirects = int(c.Proto.GetMaxRedirects())
	UniversalOptions.ReadOnly = c.Proto.GetReadOnly()
	UniversalOptions.RouteByLatency = c.Proto.GetRouteByLatency()
	UniversalOptions.RouteRandomly = c.Proto.GetRouteRandomly()
	UniversalOptions.MasterName = c.Proto.GetMasterName()
	return redis.NewUniversalClient(&UniversalOptions), nil
}
