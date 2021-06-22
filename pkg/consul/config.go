// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consul

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	net_ "github.com/searKing/golang/go/net"
	strings_ "github.com/searKing/golang/go/strings"
	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
	"github.com/searKing/sole/pkg/protobuf"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	GetViper func() *viper.Viper // If set, overrides params below
	Proto    Consul
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
		Proto: Consul{
			Address:         "",
			DefaultAddress:  "127.0.0.1:8500",
			ServiceRegistry: nil,
			ServiceResolver: nil,
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

// Validate checks Config and return a slice of found errs.
func (c *Config) Validate(validate *validator.Validate) []error {
	var errs []error
	if validate == nil {
		validate = validator.New()
	}
	errs = append(errs, validate.Struct(c))
	return errs
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

	c.Proto.Address = net_.HostportOrDefault(c.Proto.GetAddress(), c.Proto.GetDefaultAddress())
	return CompletedConfig{&completedConfig{Config: c}}
}

// NewServiceRegister creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
func (c completedConfig) NewServiceRegister() (*ServiceRegister, error) {
	if c.completeError != nil {
		return nil, c.completeError
	}
	return c.installServiceRegister()
}

// NewServiceResolver creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
func (c completedConfig) NewServiceResolver() (*ServiceResolver, error) {
	if c.completeError != nil {
		return nil, c.completeError
	}
	return c.installServiceResolver()
}

func (c *Config) loadViper() error {
	var v *viper.Viper
	if c.GetViper != nil {
		v = c.GetViper()
	}

	if err := viper_.UnmarshalProtoMessageByJsonpb(v, &c.Proto); err != nil {
		logrus.WithError(err).Errorf("load consul config from viper")
		return err
	}
	return nil
}
func (c *completedConfig) installServiceRegister() (*ServiceRegister, error) {
	register := c.Proto.GetServiceRegistry()
	var services []ServiceRegistration
	healthCheckInterval := protobuf.DurationOrDefault(register.GetHealthCheckInterval(),
		10*time.Second, "health_check_interval")
	for _, service := range register.GetServices() {
		var reg ServiceRegistration
		err := reg.SetDefault().SetAddr(
			strings_.ValueOrDefault(service.GetAddress(), register.GetDefaultServiceAddress()))
		if err != nil {
			return nil, err
		}
		reg.HealthCheckUrl = service.GetHealthCheckUrl()
		reg.HealthCheckInterval = healthCheckInterval
		reg.Complete()
		services = append(services, reg)
	}
	return NewServiceRegister(c.Proto.GetAddress(), services...)
}

func (c *completedConfig) installServiceResolver() (*ServiceResolver, error) {
	resolver := c.Proto.GetServiceResolver()
	var services []ServiceQuery
	for _, service := range resolver.GetServices() {
		var reg ServiceQuery
		reg.SetDefault()
		switch service.GetResolverType() {
		case Consul_ServiceResolver_resolver_type_random:
			reg.ResolverType = ResolverTypeRandom
		case Consul_ServiceResolver_resolver_type_consist:
			reg.ResolverType = ResolverTypeConsist
		default:
			return nil, fmt.Errorf("unsupport service_resolver_type[%s]", service.GetResolverType())
		}
		reg.Complete()
		services = append(services, reg)
	}
	return NewServiceResolver(c.Proto.GetAddress(), services...), nil
}
