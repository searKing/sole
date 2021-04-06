// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consul

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type ServiceResolverConfig struct {
	ConsulAddress string
	CheckInterval time.Duration
}

type completedServiceResolverConfig struct {
	*ServiceResolverConfig
}

type CompletedServiceResolverConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedServiceResolverConfig
}

// NewServiceResolverConfig returns a ServiceResolverConfig struct with the default values
func NewServiceResolverConfig() *ServiceResolverConfig {
	return &ServiceResolverConfig{
		CheckInterval: 10 * time.Second,
	}
}

// Validate checks ServiceResolverConfig and return a slice of found errs.
func (c *ServiceResolverConfig) Validate(validate *validator.Validate) []error {
	var errs []error
	if validate == nil {
		validate = validator.New()
	}
	errs = append(errs, validate.Struct(c))
	return errs
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
func (c *ServiceResolverConfig) Complete() CompletedServiceResolverConfig {
	var options completedServiceResolverConfig

	// set defaults
	options.ServiceResolverConfig = c
	return CompletedServiceResolverConfig{&completedServiceResolverConfig{c}}
}

func (c completedServiceResolverConfig) New() *ServiceResolver {
	resolver := NewServiceResolver(c.ConsulAddress)
	resolver.CheckInterval = c.CheckInterval
	return resolver
}
