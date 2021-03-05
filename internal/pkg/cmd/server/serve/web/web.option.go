// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/searKing/sole/internal/pkg/version"
	"github.com/searKing/sole/pkg/consul"

	"github.com/searKing/sole/internal/pkg/provider"
	"github.com/searKing/sole/pkg/webserver"
)

// ServerRunOptions runs a kubernetes api server.
type ServerRunOptions struct {
	Provider         *provider.Provider
	WebServerOptions *webserver.Config
	ServiceRegistry  *consul.Config
}

type completedServerRunOptions struct {
	*ServerRunOptions
}

// CompletedServerRunOptions is a private wrapper that enforces a call of Complete() before Run can be invoked.
type CompletedServerRunOptions struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedServerRunOptions
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

// Complete set default ServerRunOptions.
func (s *ServerRunOptions) Complete() (CompletedServerRunOptions, error) {
	if err := s.completeDiscovery(); err != nil {
		return CompletedServerRunOptions{}, err
	}
	if err := s.completeWebServer(); err != nil {
		return CompletedServerRunOptions{}, err
	}
	return CompletedServerRunOptions{&completedServerRunOptions{s}}, nil
}

// Run runs the specified APIServer.  This should never exit.
func (s *CompletedServerRunOptions) Run(ctx context.Context) error {
	// To help debugging, immediately log version
	logrus.Infof("Version: %+v", version.GetVersion())
	//isDSNAllowedOrDie(completeOptions.Provider.Proto().GetDatabase().GetDsn())

	server, err := s.WebServerOptions.Complete().New("sole")
	if err != nil {
		return err
	}

	prepared, err := server.PrepareRun()
	if err != nil {
		return err
	}

	return prepared.Run(ctx)
}
