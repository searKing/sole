// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"crypto/tls"
	"fmt"
	"net"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	gin_ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway/v2/grpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

	"github.com/searKing/sole/pkg/consul"

	"github.com/searKing/sole/pkg/net/cors"
	"github.com/searKing/sole/pkg/webserver/healthz"
)

type WebHandler interface {
	SetRoutes(ginRouter gin.IRouter, grpcRouter *grpc.Gateway)
}

type Config struct {
	GatewayOptions []grpc.GatewayOption
	GinMiddlewares []gin.HandlerFunc

	CORS *cors.Config

	TlsConfig *tls.Config

	ServiceRegistryBackend *consul.ServiceRegistryServer

	WebHandlers []WebHandler

	// done values in this values for this map are ignored.
	PostStartHooks   map[string]postStartHookEntry
	PreShutdownHooks map[string]preShutdownHookEntry

	// BindAddress is the host name to use for bind (local internet) facing URLs (e.g. Loopback)
	// Will default to a value based on secure serving info and available ipv4 IPs.
	BindAddress string
	// ExternalAddress is the host name to use for external (public internet) facing URLs (e.g. Swagger)
	// Will default to a value based on secure serving info and available ipv4 IPs.
	ExternalAddress string
	//===========================================================================
	// Fields you probably don't care about changing
	//===========================================================================

	// The default set of healthz checks. There might be more added via AddHealthChecks dynamically.
	HealthzChecks []healthz.HealthCheck
	// The default set of livez checks. There might be more added via AddHealthChecks dynamically.
	LivezChecks []healthz.HealthCheck
	// The default set of readyz-only checks. There might be more added via AddReadyzChecks dynamically.
	ReadyzChecks []healthz.HealthCheck
	// ShutdownDelayDuration allows to block shutdown for some time, e.g. until endpoints pointing to this API server
	// have converged on all node. During this time, the API server keeps serving, /healthz will return 200,
	// but /readyz will return failure.
	ShutdownDelayDuration time.Duration
}

type completedConfig struct {
	*Config
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// AddWebHandler adds a grpc and/or gin handler to our config to be exposed by the grpc gateway endpoints
// of our configured webserver.
func (c *Config) AddWebHandler(handlers ...WebHandler) {
	c.WebHandlers = append(c.WebHandlers, handlers...)
}

// AddHealthChecks adds a health check to our config to be exposed by the health endpoints
// of our configured webserver. We should prefer this to adding healthChecks directly to
// the config unless we explicitly want to add a healthcheck only to a specific health endpoint.
func (c *Config) AddHealthChecks(healthChecks ...healthz.HealthCheck) {
	for _, check := range healthChecks {
		c.HealthzChecks = append(c.HealthzChecks, check)
		c.LivezChecks = append(c.LivezChecks, check)
		c.ReadyzChecks = append(c.ReadyzChecks, check)
	}
}

// AddPostStartHook allows you to add a PostStartHook that will later be added to the server itself in a New call.
// Name conflicts will cause an error.
func (c *Config) AddPostStartHook(name string, hook PostStartHookFunc) error {
	if len(name) == 0 {
		return fmt.Errorf("missing name")
	}
	if hook == nil {
		return fmt.Errorf("hook func may not be nil: %q", name)
	}

	if postStartHook, exists := c.PostStartHooks[name]; exists {
		// this is programmer error, but it can be hard to debug
		return fmt.Errorf("unable to add %q because it was already registered by: %s", name, postStartHook.originatingStack)
	}
	c.PostStartHooks[name] = postStartHookEntry{hook: hook, originatingStack: string(debug.Stack())}

	return nil
}

// AddPostStartHookOrDie allows you to add a PostStartHook, but dies on failure.
func (c *Config) AddPostStartHookOrDie(name string, hook PostStartHookFunc) {
	if err := c.AddPostStartHook(name, hook); err != nil {
		logrus.WithError(err).Fatalf("Error registering PostStartHook %q", name)
	}
}

// AddPreShutdownHook allows you to add a PreShutdownHook.
func (s *Config) AddPreShutdownHook(name string, hook PreShutdownHookFunc) error {
	if len(name) == 0 {
		return fmt.Errorf("missing name")
	}
	if hook == nil {
		return nil
	}
	if s.PreShutdownHooks == nil {
		s.PreShutdownHooks = make(map[string]preShutdownHookEntry)
	}

	if _, exists := s.PreShutdownHooks[name]; exists {
		return fmt.Errorf("unable to add %q because it is already registered", name)
	}

	s.PreShutdownHooks[name] = preShutdownHookEntry{hook: hook}

	return nil
}

// AddPreShutdownHookOrDie allows you to add a PostStartHook, but dies on failure
func (s *Config) AddPreShutdownHookOrDie(name string, hook PreShutdownHookFunc) {
	if err := s.AddPreShutdownHook(name, hook); err != nil {
		logrus.Fatalf("Error registering PreShutdownHook %q: %v", name, err)
	}
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
func (c *Config) Complete() CompletedConfig {
	// if there is no port, and we listen on one securely, use that one
	if _, _, err := net.SplitHostPort(c.ExternalAddress); err != nil {
		if c.BindAddress == "" {
			logrus.WithError(err).Fatalf("cannot derive external address port without listening on a secure port.")
		}

		_, port, err := net.SplitHostPort(c.BindAddress)
		if err != nil {
			logrus.WithError(err).Fatalf("cannot derive external address from the secure port: %v", err)
		}
		c.ExternalAddress = net.JoinHostPort(c.ExternalAddress, port)
	}

	return CompletedConfig{&completedConfig{c}}
}

// New creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
func (c completedConfig) New(name string) (*WebServer, error) {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(
		logrus.StandardLogger().WriterLevel(logrus.InfoLevel),
		logrus.StandardLogger().WriterLevel(logrus.WarnLevel),
		logrus.StandardLogger().WriterLevel(logrus.ErrorLevel)))
	opts := grpc.WithDefault()
	opts = append(opts, grpc.WithLogrusLogger(logrus.StandardLogger()))
	opts = append(opts, grpc.WithGrpcUnaryServerChain(grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		logrus.WithError(status.Errorf(codes.Internal, "%s at %s", p, debug.Stack())).Errorf("recovered in grpc")
		return status.Errorf(codes.Internal, "%s", p)
	}))))
	opts = append(opts, grpc.WithGrpcStreamServerChain(grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		logrus.WithError(status.Errorf(codes.Internal, "%s at %s", p, debug.Stack())).Errorf("recovered in grpc")
		return status.Errorf(codes.Internal, "%s", p)
	}))))

	{
		if c.CORS != nil {
			corsHandler, err := c.CORS.Complete().New()
			if err != nil {
				return nil, err
			}
			opts = append(opts, grpc.WithHttpWrapper(corsHandler))
		}
	}

	opts = append(opts, c.GatewayOptions...)
	grpcBackend := grpc.NewGatewayTLS(c.BindAddress, c.TlsConfig, opts...)
	grpcBackend.ApplyOptions()
	ginBackend := gin.New()
	ginBackend.Use(gin.LoggerWithWriter(logrus.StandardLogger().Writer()))
	ginBackend.Use(gin_.RecoveryWithWriter(logrus.StandardLogger().Writer()))
	ginBackend.Use(gin_.UseHTTPPreflight())
	ginBackend.Use(c.GinMiddlewares...)

	s := &WebServer{
		Name:                   name,
		ServiceRegistryBackend: c.ServiceRegistryBackend,
		ShutdownDelayDuration:  c.ShutdownDelayDuration,
		grpcBackend:            grpcBackend,
		ginBackend:             ginBackend,

		postStartHooks:   map[string]postStartHookEntry{},
		preShutdownHooks: map[string]preShutdownHookEntry{},

		healthzChecks:   c.HealthzChecks,
		livezChecks:     c.LivezChecks,
		readyzChecks:    c.ReadyzChecks,
		readinessStopCh: make(chan struct{}),
	}

	// add poststarthooks that were preconfigured.  Using the add method will give us an error if the same name has already been registered.
	for name, preconfiguredPostStartHook := range c.PostStartHooks {
		if err := s.AddPostStartHook(name, preconfiguredPostStartHook.hook); err != nil {
			return nil, err
		}
	}
	for name, preconfiguredPreShutdownHook := range c.PreShutdownHooks {
		if err := s.AddPreShutdownHook(name, preconfiguredPreShutdownHook.hook); err != nil {
			return nil, err
		}
	}

	// install grpc & http handlers
	installWebHandlers(s, c.Config)

	return s, nil
}

func installWebHandlers(s *WebServer, c *Config) {
	for _, h := range c.WebHandlers {
		if h == nil {
			continue
		}
		h.SetRoutes(s.ginBackend, s.grpcBackend)
	}
}

// NewConfig returns a Config struct with the default values
func NewConfig() *Config {
	defaultHealthChecks := []healthz.HealthCheck{healthz.PingHealthCheck, healthz.LogHealthCheck}

	return &Config{
		PreShutdownHooks:      map[string]preShutdownHookEntry{},
		PostStartHooks:        map[string]postStartHookEntry{},
		HealthzChecks:         append([]healthz.HealthCheck{}, defaultHealthChecks...),
		ReadyzChecks:          append([]healthz.HealthCheck{}, defaultHealthChecks...),
		LivezChecks:           append([]healthz.HealthCheck{}, defaultHealthChecks...),
		ShutdownDelayDuration: time.Duration(0),
	}
}
