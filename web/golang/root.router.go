// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golang

import (
	"github.com/gin-gonic/gin"

	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway/v2/grpc"

	"github.com/searKing/sole/pkg/modules/opentrace"
	"github.com/searKing/sole/pkg/modules/prometheus"
	"github.com/searKing/sole/web/golang/app/configs/values"
	"github.com/searKing/sole/web/golang/app/modules/date"
	"github.com/searKing/sole/web/golang/app/modules/debug"
	"github.com/searKing/sole/web/golang/app/modules/doc/swagger"
	"github.com/searKing/sole/web/golang/app/modules/health"
	"github.com/searKing/sole/web/golang/app/modules/index"
	"github.com/searKing/sole/web/golang/app/modules/proxy"
	"github.com/searKing/sole/web/golang/app/modules/webapp"
	"github.com/searKing/sole/web/golang/app/shared/middlewares"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// SetRoutes registers this handler's routes.
func (h *Handler) SetRoutes(ginRouter gin.IRouter, grpcRouter *grpc.Gateway) {
	// bind grpcGateway as default

	ginRouter.Use(prometheus.GinHttpMetric(values.HealthMetricsPrometheusPath))
	ginRouter.Use(opentrace.GinHttpTrace(values.HealthMetricsPrometheusPath))

	middlewares.MiddlewaresRouter(ginRouter)
	index.SetRouter(ginRouter)
	debug.SetRouter(ginRouter, "")
	health.SetRouter(ginRouter)
	// webapp static files
	webapp.SetRouter(ginRouter)
	// doc
	swagger.SetRouter(ginRouter)
	// API
	apiRouter := ginRouter.Group(values.APIPathPrefix)
	index.SetRouter(apiRouter)
	debug.SetRouter(apiRouter, values.APIPathPrefix)
	health.SetRouter(apiRouter)

	date.SetRouter(grpcRouter)

	proxy.SetRouter(ginRouter)

	//// NOTE: It might be required to set SetRouter.HandleMethodNotAllowed to false to avoid problems.
	//r.HandleMethodNotAllowed = false
	//r.NotFound = Routes(h.c, values.PathPrefix)
}
