// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consul

type ServiceRegistryConfig struct {
	ConsulAddress  string // consul server addr
	ServiceAddress string
	ServiceName    string // service name
	HealthCheckUrl string
}

type completedConfig struct {
	*ServiceRegistryConfig
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// NewServiceRegistryConfig returns a Config struct with the default values
func NewServiceRegistryConfig() *ServiceRegistryConfig {
	return &ServiceRegistryConfig{
		ConsulAddress: "127.0.0.1:8500",
	}
}

// Validate checks Config and return a slice of found errs.
func (c *ServiceRegistryConfig) Validate() []error {
	return nil
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
func (c *ServiceRegistryConfig) Complete() CompletedConfig {
	var options completedConfig

	// set defaults
	options.ServiceRegistryConfig = c
	return CompletedConfig{&completedConfig{c}}
}

func (c completedConfig) New() (*ServiceRegistry, error) {
	return installConsul(c.ServiceRegistryConfig)
}

func installConsul(c *ServiceRegistryConfig) (*ServiceRegistry, error) {
	return NewServiceRegistry(c.ConsulAddress, c.ServiceName, c.ServiceAddress, c.HealthCheckUrl)
}
