// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ory/graceful"
	"github.com/searKing/sole/internal/pkg/net/serve"
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

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func ServeHTTP(handler http.Handler, address string, tlsConfig *tls.Config) (graceful.StartFunc, graceful.ShutdownFunc) {
	var srv = graceful.WithDefaults(&http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		IdleTimeout:  DefaultIdleTimeout,
		TLSConfig:    tlsConfig,
	})

	return serve.Serve(srv.Addr, srv)
}

func setupHTTP() *gin.Engine {
	return gin.New()
}
