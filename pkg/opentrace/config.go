// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opentrace

import (
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	jeagerConf "github.com/uber/jaeger-client-go/config"

	"github.com/uber/jaeger-client-go/zipkin"

	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
)

type Config struct {
	GetViper func() *viper.Viper // If set, overrides params below
	Tracing
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
		Tracing: Tracing{
			Enable: false,
			Type:   Tracing_urber_jaeger,
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
	if c.Validator == nil {
		c.Validator = validator.New()
	}
	return CompletedConfig{&completedConfig{Config: c}}
}

func (c completedConfig) Apply() (io.Closer, error) {
	if c.completeError != nil {
		return nil, c.completeError
	}
	err := c.Validate()
	if err != nil {
		return nil, err
	}
	return c.install()
}

func (c *Config) loadViper() error {
	var v *viper.Viper
	if c.GetViper != nil {
		v = c.GetViper()
	}

	if err := viper_.UnmarshalProtoMessageByJsonpb(v, &c.Tracing); err != nil {
		logrus.WithError(err).Errorf("load opentrace config from viper")
		return err
	}
	return nil
}

func (c *completedConfig) install() (io.Closer, error) {
	var Configuration jeagerConf.Configuration
	trace := c.Tracing
	if !trace.GetEnable() {
		Configuration.Disabled = true
		return Configuration.InitGlobalTracer(trace.GetServiceName())
	}

	reporter := trace.GetJaeger().GetReporter()
	if reporter != nil {
		Configuration.Reporter = &jeagerConf.ReporterConfig{}
		Configuration.Reporter.LocalAgentHostPort = reporter.GetLocalAgentHostPort()
	}

	sampler := trace.GetJaeger().GetSampler()
	if sampler != nil {
		Configuration.Sampler = &jeagerConf.SamplerConfig{}
		Configuration.Sampler.SamplingServerURL = sampler.GetServerUrl()
		Configuration.Sampler.Type = sampler.GetType().String()
		Configuration.Sampler.Param = float64(sampler.GetParam())
	}
	var configs []jeagerConf.Option
	switch trace.GetType() {
	case Tracing_urber_jaeger:
		break
	case Tracing_zipkin:
		zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
		configs = append(
			configs,
			jeagerConf.Injector(opentracing.HTTPHeaders, zipkinPropagator),
			jeagerConf.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		)
	default:
		logrus.Errorf("malformed trace type: %s", trace.GetType())
		return nil, fmt.Errorf("malformed trace type: %s", trace.GetType())
	}

	return Configuration.InitGlobalTracer(c.ServiceName, configs...)
}
