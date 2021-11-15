// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"context"

	"github.com/go-playground/validator/v10"
	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
	viperpb "github.com/searKing/sole/api/protobuf-spec/v1/viper"
	"github.com/spf13/viper"
)

//go:generate go-option -type=Config
type Config struct {
	Proto     *viperpb.ViperProto
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
	return &Config{
		Proto: &viperpb.ViperProto{},
	}
}

// NewViperConfig returns a Config struct with the viper instance
// key representing a subtree of this instance.
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
	return CompletedConfig{&completedConfig{Config: c}}
}

// New creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
// New usually called after Complete
func (c completedConfig) New(ctx context.Context) (*Provider, error) {
	p := &Provider{
		Proto: c.Proto,
		ctx:   ctx,
	}
	return p, nil
}
