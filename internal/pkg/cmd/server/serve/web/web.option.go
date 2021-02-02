// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"github.com/searKing/sole/pkg/modules/consul"

	"github.com/searKing/sole/internal/pkg/provider"
	"github.com/searKing/sole/pkg/modules/webserver"
)

// ServerRunOptions runs a kubernetes api server.
type ServerRunOptions struct {
	Provider         *provider.Provider
	WebServerOptions *webserver.Config
	ServiceRegistry  *consul.Config
}

func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{
		Provider:         provider.GlobalProvider(),
		WebServerOptions: webserver.NewConfig(),
		ServiceRegistry:  consul.NewConfig(),
	}
}

// Validate checks ServerRunOptions and return a slice of found errs.
func (s *ServerRunOptions) Validate() []error {
	var errs []error
	errs = append(errs, s.ServiceRegistry.Validate()...)
	return errs
}

// CompletedServerRunOptions is a private wrapper that enforces a call of Complete() before Run can be invoked.
type CompletedServerRunOptions struct {
	*ServerRunOptions
}

// Complete set default ServerRunOptions.
func Complete(s *ServerRunOptions) (CompletedServerRunOptions, error) {
	var options CompletedServerRunOptions
	if err := s.completeDiscovery(); err != nil {
		return options, err
	}
	if err := s.completeWebServer(); err != nil {
		return options, err
	}

	options.ServerRunOptions = s
	return options, nil
}
