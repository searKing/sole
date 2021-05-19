// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"context"
	"crypto/tls"
	"net"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/searKing/golang/go/time/rate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	grpc_ "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

	gin_ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway/v2/grpc"
	logrus_ "github.com/searKing/golang/third_party/github.com/sirupsen/logrus"
	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"

	"github.com/searKing/sole/pkg/protobuf"

	"github.com/searKing/sole/pkg/consul"

	"github.com/searKing/sole/pkg/net/cors"
	"github.com/searKing/sole/pkg/webserver/healthz"
)

type Config struct {
	GetViper func() *viper.Viper // If set, overrides params below
	Proto    Web

	GatewayOptions []grpc.GatewayOption
	GinMiddlewares []gin.HandlerFunc

	CORS *cors.Config

	TlsConfig *tls.Config

	ServiceRegistryBackend *consul.ServiceRegister
	ServiceResolverBackend *consul.ServiceResolver

	// BindAddress is the host name to use for bind (local internet) facing URLs (e.g. Loopback)
	// Will default to a value based on secure serving info and available ipv4 IPs.
	BindAddress string
	// ExternalAddress is the host name to use for external (public internet) facing URLs (e.g. Swagger)
	// Will default to a value based on secure serving info and available ipv4 IPs.
	ExternalAddress string
	// ShutdownDelayDuration allows to block shutdown for some time, e.g. until endpoints pointing to this API server
	// have converged on all node. During this time, the API server keeps serving, /healthz will return 200,
	// but /readyz will return failure.
	ShutdownDelayDuration time.Duration
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

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
func (c *Config) Complete() CompletedConfig {
	if err := c.loadViper(); err != nil {
		return CompletedConfig{&completedConfig{
			Config:        c,
			completeError: err,
		}}
	}
	c.parseViper()
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

	return CompletedConfig{&completedConfig{Config: c}}
}

// New creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
func (c completedConfig) New(name string) (*WebServer, error) {
	if c.completeError != nil {
		return nil, c.completeError
	}
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(
		logrus.StandardLogger().WriterLevel(logrus.InfoLevel),
		logrus.StandardLogger().WriterLevel(logrus.WarnLevel),
		logrus.StandardLogger().WriterLevel(logrus.ErrorLevel)))
	opts := grpc.WithDefault()
	if c.Proto.GetNoGrpcProxy() {
		opts = append(opts, grpc.WithGrpcDialOption(grpc_.WithNoProxy()))
	}
	opts = append(opts, grpc.WithLogrusLogger(logrus.StandardLogger()))
	opts = append(opts, grpc.WithGrpcUnaryServerChain(grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		logrus.WithError(status.Errorf(codes.Internal, "%s at %s", p, debug.Stack())).Errorf("recovered in grpc")
		return status.Errorf(codes.Internal, "%s", p)
	}))))
	opts = append(opts, grpc.WithGrpcStreamServerChain(grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		logrus.WithError(status.Errorf(codes.Internal, "%s at %s", p, debug.Stack())).Errorf("recovered in grpc")
		return status.Errorf(codes.Internal, "%s", p)
	}))))

	if c.Proto.GetMaxConcurrencyUnary() > 0 {
		limiter := rate.NewFullBurstLimiter(int(c.Proto.GetMaxConcurrencyUnary()))
		opts = append(opts, grpc.WithGrpcUnaryServerChain(
			func(ctx context.Context, req interface{},
				info *grpc_.UnaryServerInfo, handler grpc_.UnaryHandler) (resp interface{}, err error) {
				if limiter.Allow() {
					defer limiter.PutToken()
					return handler(ctx, req)
				}
				err = status.Errorf(codes.ResourceExhausted,
					"%s is rejected by ratelimit middleware, please retry later", info.FullMethod)
				logrus.WithError(err).Errorf("refuesd by grpc")
				return nil, err
			}))
	}
	if c.Proto.GetMaxConcurrencyStream() > 0 {
		limiter := rate.NewFullBurstLimiter(int(c.Proto.GetMaxConcurrencyStream()))
		opts = append(opts, grpc.WithGrpcStreamServerChain(
			func(srv interface{}, ss grpc_.ServerStream, info *grpc_.StreamServerInfo, handler grpc_.StreamHandler) error {
				if limiter.Allow() {
					defer limiter.PutToken()
					return handler(srv, ss)
				}
				err := status.Errorf(codes.ResourceExhausted,
					"%s is rejected by ratelimit middleware, please retry later", info.FullMethod)
				logrus.WithError(err).Errorf("refuesd by grpc")
				return err
			}))
	}
	{
		if c.CORS != nil {
			opts = append(opts, grpc.WithHttpWrapper(c.CORS.Complete().New().Handler))
			c.GinMiddlewares = append(c.GinMiddlewares, c.CORS.Complete().NewGinHandler())
		}
	}

	opts = append(opts, c.GatewayOptions...)
	grpcBackend := grpc.NewGatewayTLS(c.BindAddress, c.TlsConfig, opts...)
	grpcBackend.ApplyOptions()
	grpcBackend.ErrorLog = logrus_.AsStdLogger(logrus.StandardLogger(), logrus.ErrorLevel, "", 0)
	ginBackend := gin.New()
	ginBackend.Use(gin.LoggerWithWriter(logrus.StandardLogger().Writer()))
	ginBackend.Use(gin_.RecoveryWithWriter(grpcBackend.ErrorLog.Writer()))
	ginBackend.Use(gin_.UseHTTPPreflight())
	ginBackend.Use(c.GinMiddlewares...)

	defaultHealthChecks := []healthz.HealthCheck{healthz.PingHealthCheck, healthz.LogHealthCheck}

	s := &WebServer{
		Name:                  name,
		ShutdownDelayDuration: c.ShutdownDelayDuration,
		grpcBackend:           grpcBackend,
		ginBackend:            ginBackend,

		postStartHooks:   map[string]postStartHookEntry{},
		preShutdownHooks: map[string]preShutdownHookEntry{},
		healthzChecks:    defaultHealthChecks,
		livezChecks:      defaultHealthChecks,
		readyzChecks:     defaultHealthChecks,
		readinessStopCh:  make(chan struct{}),
	}

	return s, nil
}

// NewConfig returns a Config struct with the default values
func NewConfig() *Config {

	return &Config{
		ShutdownDelayDuration: time.Duration(0),
		Proto: Web{
			BindAddr: &Web_Net{
				Port: 80,
			},
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

func (c *Config) loadViper() error {
	var v *viper.Viper
	if c.GetViper != nil {
		v = c.GetViper()
	}

	if err := viper_.UnmarshalProtoMessageByJsonpb(v, &c.Proto); err != nil {
		logrus.WithError(err).Errorf("load web_server config from viper")
		return err
	}
	return nil
}

func (s *Config) parseViper() {
	s.BindAddress = s.Proto.GetBackendBindHostPort()
	s.ExternalAddress = s.Proto.GetBackendServeHostPort(true)

	{
		corsConfig := cors.NewConfig()
		corsInfo := s.Proto.GetCors()
		if corsInfo != nil {
			if corsInfo.Enable {
				maxAge := protobuf.DurationOrDefault(corsInfo.GetMaxAge(), 0, "max_age")
				corsConfig.UseConditional = corsInfo.GetUseConditional()
				corsConfig.AllowedOrigins = corsInfo.GetAllowedOrigins()
				corsConfig.AllowedMethods = corsInfo.GetAllowedOrigins()
				corsConfig.AllowedHeaders = corsInfo.GetAllowedHeaders()
				corsConfig.ExposedHeaders = corsInfo.GetExposedHeaders()

				corsConfig.MaxAge = maxAge
				corsConfig.AllowCredentials = corsInfo.GetAllowCredentials()
			} else {
				corsConfig = nil
			}
		}
		s.CORS = corsConfig
	}
}
