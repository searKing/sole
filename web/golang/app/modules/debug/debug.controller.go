// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debug

import (
	_ "expvar"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/searKing/golang/go/version"
	gin2 "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
	"github.com/searKing/sole/web/golang/app/configs/values"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	pathPrefixTrim string
}

func NewController() *Controller {
	return &Controller{}
}

// SetRoutes registers this handler's routes.
func (c *Controller) SetRoutes(ginRouter gin.IRouter, grpcRouter *grpc.Gateway) {
	logrus.Info("installing router")

	ginRouter.GET(values.DebugPProf, c.PProf())
	ginRouter.GET(values.DebugExpVar, c.ExpVar())
	ginRouter.GET(values.DebugMetricsPrometheusPath, c.MetricsPrometheus())
	ginRouter.GET(values.DebugVersionPath, c.Version())
}

func (c *Controller) PProf() gin.HandlerFunc {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	if c.pathPrefixTrim != "" {
		return gin2.RedirectTrim(http.StatusFound, c.pathPrefixTrim)
	}
	return gin.WrapH(http.DefaultServeMux)
}

func (c *Controller) ExpVar() gin.HandlerFunc {
	if c.pathPrefixTrim != "" {
		return gin2.RedirectTrim(http.StatusFound, c.pathPrefixTrim)
	}
	return gin.WrapH(http.DefaultServeMux)
}

// MetricsPrometheus Prometheus指标统计
func (c *Controller) MetricsPrometheus() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func (c *Controller) Version() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "%s", version.Get().String())
	}
}
