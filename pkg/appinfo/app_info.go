// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package appinfo

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	GetViper  func() *viper.Viper // If set, overrides params below
	Proto     AppInfo
	Validator *validator.Validate
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
		Proto: AppInfo{
			Version:            Version,
			GitHash:            GitHash,
			BuildTime:          BuildTime,
			ServiceName:        ServiceName,
			ServiceDisplayName: ServiceDisplayName,
			ServiceDescription: ServiceDescription,
			ServiceId:          ServiceId,
		},
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

// Validate checks Config.
func (c *completedConfig) Validate() error {
	return c.Validator.Struct(c)
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
func (c *Config) Complete() CompletedConfig {
	if err := c.loadViper(); err != nil {
		return CompletedConfig{&completedConfig{
			Config:        c,
			completeError: err,
		}}
	}

	if c.Proto.GetServiceId() == "" {
		c.Proto.ServiceId = uuid.New().String()
	}
	if c.Validator == nil {
		c.Validator = validator.New()
	}
	return CompletedConfig{&completedConfig{Config: c}}
}

// Apply creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
func (c completedConfig) Apply() error {
	if c.completeError != nil {
		return c.completeError
	}
	err := c.Validate()
	if err != nil {
		return err
	}
	return c.install()
}

func (c *Config) loadViper() error {
	var v *viper.Viper
	if c.GetViper != nil {
		v = c.GetViper()
	}

	if err := viper_.UnmarshalProtoMessageByJsonpb(v, &c.Proto); err != nil {
		logrus.WithError(err).Errorf("load logs config from viper")
		return err
	}
	return nil
}

func (c *completedConfig) install() error {
	Version = c.Proto.GetVersion()
	BuildTime = c.Proto.GetBuildTime()
	GitHash = c.Proto.GetGitHash()

	ServiceName = c.Proto.GetServiceName()
	ServiceDisplayName = c.Proto.GetServiceDisplayName()
	ServiceDescription = c.Proto.GetServiceDescription()
	ServiceId = c.Proto.GetServiceId()
	return nil
}