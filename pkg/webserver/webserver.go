// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/searKing/golang/go/x/graceful"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway/v2/grpc"
	"github.com/sirupsen/logrus"

	"github.com/searKing/sole/pkg/webserver/healthz"
)

type WebHandler interface {
	SetRoutes(ginRouter gin.IRouter, grpcRouter *grpc.Gateway)
}

type WebServer struct {
	Name string

	ginBackend  *gin.Engine
	grpcBackend *grpc.Gateway

	// PostStartHooks are each called after the server has started listening, in a separate go func for each
	// with no guarantee of ordering between them.  The map key is a name used for error reporting.
	// It may kill the process with a panic if it wishes to by returning an error.
	postStartHookLock    sync.Mutex
	postStartHooks       map[string]postStartHookEntry
	postStartHooksCalled bool

	preShutdownHookLock    sync.Mutex
	preShutdownHooks       map[string]preShutdownHookEntry
	preShutdownHooksCalled bool

	// healthz checks
	healthzLock            sync.Mutex
	healthzChecks          []healthz.HealthCheck
	healthzChecksInstalled bool
	// livez checks
	livezLock            sync.Mutex
	livezChecks          []healthz.HealthCheck
	livezChecksInstalled bool
	// readyz checks
	readyzLock            sync.Mutex
	readyzChecks          []healthz.HealthCheck
	readyzChecksInstalled bool
	livezGracePeriod      time.Duration

	// the readiness stop channel is used to signal that the apiserver has initiated a shutdown sequence, this
	// will cause readyz to return unhealthy.
	readinessStopCh chan struct{}

	// ShutdownDelayDuration allows to block shutdown for some time, e.g. until endpoints pointing to this API server
	// have converged on all node. During this time, the API server keeps serving, /healthz will return 200,
	// but /readyz will return failure.
	ShutdownDelayDuration time.Duration
}

func NewWebServer(ctx context.Context, config *Config) (*WebServer, error) {
	return config.Complete().New()
}

// preparedWebServer is a private wrapper that enforces a call of PrepareRun() before Run can be invoked.
type preparedWebServer struct {
	*WebServer
}

// PrepareRun does post API installation setup steps. It calls recursively the same function of the delegates.
func (s *WebServer) PrepareRun() (preparedWebServer, error) {
	if s.grpcBackend != nil {
		s.grpcBackend.Handler = s.ginBackend
	}

	s.installHealthz()
	s.installLivez()
	err := s.addReadyzShutdownCheck(s.readinessStopCh)
	if err != nil {
		logrus.Errorf("Failed to parseViper readyz shutdown check %s", err)
		return preparedWebServer{}, err
	}
	s.installReadyz()

	// Register audit backend preShutdownHook.
	return preparedWebServer{s}, nil
}

// Run spawns the secure http server. It only returns if stopCh is closed
// or the secure port cannot be listened on initially.
func (s preparedWebServer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()

		<-ctx.Done()

		// As soon as shutdown is initiated, /readyz should start returning failure.
		// This gives the load balancer a window defined by ShutdownDelayDuration to detect that /readyz is red
		// and stop sending traffic to this server.
		close(s.readinessStopCh)

		time.Sleep(s.ShutdownDelayDuration)
	}()

	// close socket after delayed stopCh
	ctx, err := s.NonBlockingRun(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()

	// run shutdown hooks directly. This includes deregistering from the kubernetes endpoint in case of kube-apiserver.
	err = s.RunPreShutdownHooks()
	if err != nil {
		return err
	}

	// wait for the delayed stopCh before closing the handler chain (it rejects everything after Wait has been called).
	wg.Wait()

	return nil
}

// NonBlockingRun spawns the secure http|grpc server. An error is
// returned if the secure port cannot be listened on.
// The returned context is done when the (asynchronous) termination is finished.
func (s preparedWebServer) NonBlockingRun(ctx context.Context) (context.Context, error) {
	ctx, cancel := context.WithCancel(ctx)

	// Start the audit backend before any request comes in. This means we must call Backend.Run
	// before http server start serving. Otherwise the Backend.ProcessEvents call might block.

	go func() {
		defer cancel()
		err := graceful.Graceful(ctx, graceful.Handler{
			Name: s.Name,
			StartFunc: func(ctx context.Context) error {
				var err error
				logrus.Infof("Setting up http server on %s", s.grpcBackend.Addr)
				err = s.grpcBackend.ListenAndServe()
				if err != nil {
					logrus.WithError(err).Errorf("Have not set up http server on %s", s.grpcBackend.Addr)
					return err
				}
				return nil
			},
			ShutdownFunc: func(ctx context.Context) error {
				logrus.Infof("Shutting down http server on %s", s.grpcBackend.Addr)
				err := s.grpcBackend.Shutdown(ctx)
				if err != nil {
					logrus.WithError(err).Errorf("Have not shut down http server on %s", s.grpcBackend.Addr)
					return err
				}
				logrus.Infof("Have Shut down http server on %s", s.grpcBackend.Addr)
				return nil
			},
		})
		if err != nil {
			logrus.WithError(err).Fatal("Could not gracefully run servers")
		}
	}()

	s.RunPostStartHooks(ctx)

	return ctx, nil
}

func (s *WebServer) InstallWebHandlers(handlers ...WebHandler) {
	for _, h := range handlers {
		if h == nil {
			continue
		}
		h.SetRoutes(s.ginBackend, s.grpcBackend)
	}
}
