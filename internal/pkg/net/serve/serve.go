// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serve

import (
	"context"
	"net"
	"time"

	"github.com/ory/graceful"
	"github.com/searKing/sole/internal/pkg/provider"
)

const (
	// DefaultReadTimeout sets the maximum time a client has to fully stream a request (30m)
	DefaultReadTimeout = 30 * time.Minute
	// DefaultWriteTimeout sets the maximum amount of time a handler has to fully process a request (1h)
	DefaultWriteTimeout = 1 * time.Hour
	// DefaultIdleTimeout sets the maximum amount of time a Keep-Alive connection can remain idle before
	// being recycled (2h)
	DefaultIdleTimeout = 2 * time.Hour
)

type Server interface {
	Serve(l net.Listener) error
	ListenAndServe() error
	ListenAndServeTLS(certFile, keyFile string) error
	Shutdown(ctx context.Context) error
	RegisterOnShutdown(f func())
}

func Serve(addr string, srv Server) (graceful.StartFunc, graceful.ShutdownFunc) {
	srv.RegisterOnShutdown(func() {
		_ = provider.GlobalProvider().Tracer().Close()
	})

	return func() error {
			logger := provider.GlobalProvider().Logger()
			webInfo := provider.GlobalProvider().Proto().GetWeb()

			logger.Infof("Setting up http server on %s", addr)

			if webInfo.GetForceDisableTls() {
				logger.Warnln("HTTPS disabled. Never do this in production.")
				return srv.ListenAndServe()
			}
			if len(webInfo.GetTls().GetAllowedTlsCidrs()) != 0 {
				logger.Infoln("TLS termination enabled, disabling https.")
				return srv.ListenAndServe()
			}
			return srv.ListenAndServeTLS("", "")
		}, func(ctx context.Context) error {
			logger := provider.GlobalProvider().Logger()

			logger.Infof("Shutting down http server on %s", addr)
			defer logger.Infof("Have Shut down http server on %s", addr)
			return srv.Shutdown(ctx)
		}
}
