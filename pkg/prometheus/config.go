// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	Name      string
	Version   string
	Hash      string
	BuildTime string
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
	return &Config{}
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

func (c completedConfig) New() (http.Handler, error) {
	// Add Go module build info.
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
	return promhttp.Handler(), nil
}
