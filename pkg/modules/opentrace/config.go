// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opentrace

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	jeagerConf "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/zipkin"
)

type Config struct {
	Enabled       bool
	ServiceName   string
	Type          Type
	Configuration jeagerConf.Configuration

	closer io.Closer
}

type completedConfig struct {
	*Config
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// NewConfig returns a Config struct with the default values
func NewConfig() *Config {
	return &Config{
		Type:          TypeJeager,
		Configuration: jeagerConf.Configuration{},
	}
}

// Validate checks Config and return a slice of found errs.
func (s *Config) Validate() []error {
	return nil
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
func (s *Config) Complete() CompletedConfig {
	var options completedConfig

	// set defaults
	options.Config = s
	return CompletedConfig{&completedConfig{s}}
}

func (c completedConfig) Apply() (io.Closer, error) {
	return installOpenTrace(c.Config)
}

func installOpenTrace(c *Config) (io.Closer, error) {
	if !c.Enabled {
		c.Configuration.Disabled = true
		return c.Configuration.InitGlobalTracer(c.ServiceName)
	}
	if c.Type >= TypeButt {
		return nil, fmt.Errorf("unknown tracer: %s", c.Type)
	}
	var configs []jeagerConf.Option
	if c.Type == TypeZipkin {
		zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
		configs = append(
			configs,
			jeagerConf.Injector(opentracing.HTTPHeaders, zipkinPropagator),
			jeagerConf.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		)
	}

	return c.Configuration.InitGlobalTracer(c.ServiceName, configs...)
}
