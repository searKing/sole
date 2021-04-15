// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/searKing/sole/pkg/consul"
	"github.com/searKing/sole/web/golang"
	"github.com/sirupsen/logrus"

	"github.com/searKing/sole/internal/pkg/provider"
	"github.com/searKing/sole/internal/pkg/version"
	"github.com/searKing/sole/pkg/webserver"
)

// ServerRunOptions runs a kubernetes api server.
type ServerRunOptions struct {
	Provider         *provider.Provider
	WebServerOptions *webserver.Config
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
		WebServerOptions: webserver.NewViperConfig("web"),
	}
}

// Validate checks ServerRunOptions and return a slice of found errs.
func (s *ServerRunOptions) Validate(validate *validator.Validate) []error {
	var errs []error
	return errs
}

// Complete set default ServerRunOptions.
func (s *ServerRunOptions) Complete() (CompletedServerRunOptions, error) {
	s.WebServerOptions.Proto.ForceDisableTls = provider.ForceDisableTls
	return CompletedServerRunOptions{&completedServerRunOptions{s}}, nil
}

// Run runs the specified APIServer.  This should never exit.
func (s *CompletedServerRunOptions) Run(ctx context.Context) error {
	// To help debugging, immediately log version
	logrus.Infof("Version: %+v", version.GetVersion())
	//isDSNAllowedOrDie(completeOptions.Provider.Proto.GetDatabase().GetDsn())

	server, err := s.WebServerOptions.Complete().New("sole")
	if err != nil {
		return err
	}
	server.InstallWebHandlers(golang.NewHandler())
	{
		// register webserver as a service
		if register := s.Provider.ServiceRegister; register != nil {
			for _, domain := range s.WebServerOptions.Proto.GetAdvertiseAddr().GetDomains() {
				if domain == "" {
					continue
				}
				r := consul.ServiceRegistration{}
				if err := r.SetDefault().SetAddr(s.WebServerOptions.Proto.GetBackendServeHostPort()); err != nil {
					return err
				}
				r.Name = domain
				r.HealthCheckUrl = "/healthz"
				r.Complete()
				if err := register.AddService(r); err != nil {
					logrus.WithError(err).WithField("domain", domain).
						Errorf("register web server as service in consul")
					return err
				}
			}

			server.AddPostStartHookOrDie("service-register-backend", func(ctx context.Context) error {
				return register.Run(ctx)
			})
			server.AddPreShutdownHookOrDie("service-register-backend", func() error {
				register.Shutdown()
				return nil
			})
		}
	}
	{
		// resolve webserver as a service
		if resolver := s.Provider.ServiceResolver; resolver != nil {

			for _, domain := range s.WebServerOptions.Proto.GetAdvertiseAddr().GetDomains() {
				if domain == "" {
					continue
				}
				r := consul.ServiceQuery{}
				r.SetDefault()
				r.Name = domain
				r.Complete()
				if err := resolver.AddService(r); err != nil {
					logrus.WithError(err).Errorf("resolver web server as service in consul")
					return err
				}
			}

			server.AddPostStartHookOrDie("service-resolver-backend", func(ctx context.Context) error {
				return resolver.Run(ctx)
			})
			server.AddPreShutdownHookOrDie("service-resolver-backend", func() error {
				resolver.Shutdown()
				return nil
			})
		}
	}

	prepared, err := server.PrepareRun()
	if err != nil {
		return err
	}

	return prepared.Run(ctx)
}
