// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package health

import (
	"github.com/gin-gonic/gin"
	"github.com/searKing/sole/web/golang/app/configs/values"
)

func SetRouter(router gin.IRouter) gin.IRouter {
	health := NewController()
	router.GET(values.HealthAliveCheckPath, health.Alive())
	router.GET(values.HealthReadyCheckPath, health.Ready(true))
	router.GET(values.HealthMetricsPrometheusPath, health.MetricsPrometheus())

	router.GET(values.HealthVersionPath, health.Version())
	return router
}
