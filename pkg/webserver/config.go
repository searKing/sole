// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"crypto/tls"
	"fmt"
	"math"
	"net"
	"os"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/searKing/golang/third_party/google.golang.org/grpc/interceptors/burstlimit"
	"github.com/searKing/golang/third_party/google.golang.org/grpc/interceptors/timeoutlimit"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

	gin_ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
	grpc_ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
	logrus_ "github.com/searKing/golang/third_party/github.com/sirupsen/logrus"
	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"

	"github.com/searKing/sole/pkg/consul"

	"github.com/searKing/sole/pkg/net/cors"
	"github.com/searKing/sole/pkg/webserver/healthz"
)

// ClientMaxReceiveMessageSize use 4GB as the default message size limit.
// grpc library default is 4MB
var defaultMaxReceiveMessageSize = math.MaxInt32 // 1024 * 1024 * 4
var defaultMaxSendMessageSize = math.MaxInt32

type Config struct {
	Proto     Web
	viper     *viper.Viper
	viperKeys []string

	GatewayOptions []grpc_.GatewayOption
	GinMiddlewares []gin.HandlerFunc

	CORS *cors.Config

	TlsConfig *tls.Config

	ServiceRegistryBackend *consul.ServiceRegister
	ServiceResolverBackend *consul.ServiceResolver

	// Name is the human-readable server name, optional
	Name string
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
	if c.viper != nil {
		err := viper_.UnmarshalKeysViper(c.viper, c.viperKeys, &c.Proto)
		if err != nil {
			return CompletedConfig{&completedConfig{completeError: err}}
		}
	}
	c.BindAddress = c.Proto.GetBackendBindHostPort()
	c.ExternalAddress = c.Proto.GetBackendServeHostPort(true)

	{
		c.CORS = cors.NewViperConfig(c.viper, append(c.viperKeys, "cors")...)
	}

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
func (c completedConfig) New() (*WebServer, error) {
	if c.completeError != nil {
		return nil, c.completeError
	}
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(
		logrus.StandardLogger().WriterLevel(logrus.DebugLevel),
		logrus.StandardLogger().WriterLevel(logrus.WarnLevel),
		logrus.StandardLogger().WriterLevel(logrus.ErrorLevel)))
	opts := grpc_.WithDefault()
	if c.Proto.GetNoGrpcProxy() {
		opts = append(opts, grpc_.WithGrpcDialOption(grpc.WithNoProxy()))
	}
	opts = append(opts, grpc_.WithLogrusLogger(logrus.StandardLogger()))
	{
		// 设置GRPC最大消息大小
		opts = append(opts, grpc_.WithGrpcDialOption(grpc.WithNoProxy()))
		// http -> grpc client -> grpc server
		if c.Proto.GetMaxReceiveMessageSizeInBytes() > 0 {
			opts = append(opts, grpc_.WithGrpcServerOption(grpc.MaxRecvMsgSize(int(c.Proto.GetMaxReceiveMessageSizeInBytes()))))
			opts = append(opts, grpc_.WithGrpcDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(int(c.Proto.GetMaxReceiveMessageSizeInBytes())))))
		} else {
			opts = append(opts, grpc_.WithGrpcServerOption(grpc.MaxRecvMsgSize(defaultMaxReceiveMessageSize)))
			opts = append(opts, grpc_.WithGrpcDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(defaultMaxReceiveMessageSize))))
		}
		// http <- grpc client <- grpc server
		if c.Proto.GetMaxSendMessageSizeInBytes() > 0 {
			opts = append(opts, grpc_.WithGrpcServerOption(grpc.MaxSendMsgSize(int(c.Proto.GetMaxSendMessageSizeInBytes()))))
			opts = append(opts, grpc_.WithGrpcDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(c.Proto.GetMaxSendMessageSizeInBytes())))))
		} else {
			opts = append(opts, grpc_.WithGrpcServerOption(grpc.MaxSendMsgSize(defaultMaxSendMessageSize)))
			opts = append(opts, grpc_.WithGrpcDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(defaultMaxSendMessageSize))))
		}
	}
	{
		// recover
		opts = append(opts, grpc_.WithGrpcUnaryServerChain(grpcrecovery.UnaryServerInterceptor(grpcrecovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logrus.WithError(status.Errorf(codes.Internal, "%s at %s", p, debug.Stack())).Errorf("recovered in grpc")
			{
				_, _ = os.Stderr.Write([]byte(fmt.Sprintf("panic: %s", p)))
				debug.PrintStack()
				_, _ = os.Stderr.Write([]byte(" [recovered]"))
				_, _ = os.Stderr.Write([]byte("\n"))
			}
			return status.Errorf(codes.Internal, "%s", p)
		}))))
		opts = append(opts, grpc_.WithGrpcStreamServerChain(grpcrecovery.StreamServerInterceptor(grpcrecovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logrus.WithError(status.Errorf(codes.Internal, "%s at %s", p, debug.Stack())).Errorf("recovered in grpc")
			{
				_, _ = os.Stderr.Write([]byte(fmt.Sprintf("panic: %s", p)))
				debug.PrintStack()
				_, _ = os.Stderr.Write([]byte(" [recovered]"))
				_, _ = os.Stderr.Write([]byte("\n"))
			}
			return status.Errorf(codes.Internal, "%s", p)
		}))))
	}
	{
		// handle request timeout
		opts = append(opts, grpc_.WithGrpcUnaryServerChain(
			timeoutlimit.UnaryServerInterceptor(c.Proto.GetHandledTimeoutUnary().AsDuration())))
		opts = append(opts, grpc_.WithGrpcStreamServerChain(
			timeoutlimit.StreamServerInterceptor(c.Proto.GetHandledTimeoutStream().AsDuration())))
	}
	{
		// burst limit
		opts = append(opts, grpc_.WithGrpcUnaryServerChain(
			burstlimit.UnaryServerInterceptor(int(c.Proto.GetMaxConcurrencyUnary()), c.Proto.GetBurstLimitTimeoutUnary().AsDuration())))
		opts = append(opts, grpc_.WithGrpcStreamServerChain(
			burstlimit.StreamServerInterceptor(int(c.Proto.GetMaxConcurrencyStream()), c.Proto.GetBurstLimitTimeoutStream().AsDuration())))
	}
	{
		if c.CORS != nil {
			opts = append(opts, grpc_.WithHttpWrapper(c.CORS.Complete().New().Handler))
			c.GinMiddlewares = append(c.GinMiddlewares, c.CORS.Complete().NewGinHandler())
		}
	}

	opts = append(opts, c.GatewayOptions...)
	grpcBackend := grpc_.NewGatewayTLS(c.BindAddress, c.TlsConfig, opts...)
	grpcBackend.ApplyOptions()
	grpcBackend.ErrorLog = logrus_.AsStdLogger(logrus.StandardLogger(), logrus.ErrorLevel, "", 0)
	ginBackend := gin.New()
	ginBackend.Use(gin.LoggerWithWriter(logrus.StandardLogger().Writer()))
	ginBackend.Use(gin_.RecoveryWithWriter(grpcBackend.ErrorLog.Writer()))
	ginBackend.Use(gin_.UseHTTPPreflight())
	ginBackend.Use(c.GinMiddlewares...)

	defaultHealthChecks := []healthz.HealthCheck{healthz.PingHealthCheck, healthz.LogHealthCheck}

	s := &WebServer{
		Name:                  c.Name,
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
// key representing a subtree of this instance.
// NewViperConfig is case-insensitive for a key.
func NewViperConfig(v *viper.Viper, keys ...string) *Config {
	c := NewConfig()
	c.viper = v
	c.viperKeys = keys
	return c
}
