// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debug

import (
	_ "expvar" // for expvar
	"log/slog"
	"net/http"
	_ "net/http/pprof" // for pprof
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/searKing/golang/go/version"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
)

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

// SetRoutes registers this handler's routes.
func (c *Controller) SetRoutes(ginRouter gin.IRouter, grpcRouter *grpc.Gateway) {
	slog.Info("installing router")
	if EnablePprof {
		ginRouter.GET("/debug/pprof/*path", c.PProf())
	}
	ginRouter.GET("/debug/vars", c.ExpVar())
	ginRouter.Any("/metrics/prometheus/*path", c.MetricsPrometheus())
	ginRouter.Any("/version", c.Version())
}

// PProf serves via its HTTP server runtime profiling data in the format expected by the pprof visualization tool.
func (c *Controller) PProf() gin.HandlerFunc {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	return gin.WrapH(http.DefaultServeMux)
}

// ExpVar exposes public expvar variables via HTTP at /debug/vars.
func (c *Controller) ExpVar() gin.HandlerFunc {
	return gin.WrapH(http.DefaultServeMux)
}

// MetricsPrometheus Prometheus Metrics
func (c *Controller) MetricsPrometheus() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

// Version version
func (c *Controller) Version() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, version.Get().String())
	}
}
