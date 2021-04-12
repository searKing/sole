// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package health

import (
	_ "expvar"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func (d *Controller) Health() gin.HandlerFunc {
	router := httprouter.New()
	//d.h.SetRoutes(router, true)
	return gin.WrapF(router.ServeHTTP)
}

func (d *Controller) MetricsPrometheus() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func (d *Controller) Alive() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (d *Controller) Ready(shareErrors bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (d *Controller) Version() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
