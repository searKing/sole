// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package health

import (
	_ "expvar"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/x/healthx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/serverx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/searKing/golang/pkg/net/http"
	"github.com/searKing/sole/internal/pkg/provider"
)

type HealthController struct {
	h *healthx.Handler
}

func NewHealthController() *HealthController {
	h := &HealthController{}
	h.init()
	return h
}

func (d *HealthController) init() {
	c := provider.GlobalProvider()
	//dependency.ExpectDependency(c.Logger(), map[string]interface{}{"service_discovery": c.ServiceDiscoveryConnection})

	logger := logrusx.New("", "")
	logger.Logger = c.Logger()
	w := herodot.NewJSONWriter(logger)
	w.ErrorEnhancer = serverx.ErrorEnhancerRFC6749
	d.h = healthx.NewHandler(w, c.Proto().GetAppInfo().GetBuildVersion(), healthx.ReadyCheckers{
		"database": provider.GlobalProvider().SqlDBPing,
		//"zookeeper": ctx.ServiceDiscoveryConnection.Ping,
	})
}

func (d *HealthController) Health() gin.HandlerFunc {
	router := httprouter.New()
	d.h.SetRoutes(router, true)
	return gin.WrapF(router.ServeHTTP)
}

func (d *HealthController) MetricsPrometheus() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func (d *HealthController) Alive() gin.HandlerFunc {
	return gin.WrapH(http.WrapHTTPRouterF(d.h.Alive))
}

func (d *HealthController) Ready(shareErrors bool) gin.HandlerFunc {
	return gin.WrapH(http.WrapHTTPRouterF(d.h.Ready(shareErrors)))
}

func (d *HealthController) Version() gin.HandlerFunc {
	return gin.WrapH(http.WrapHTTPRouterF(d.h.Version))
}
