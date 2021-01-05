// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middlewares

import (
	"github.com/gin-gonic/gin"
	gin_ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
	"github.com/searKing/sole/internal/pkg/provider"
	"github.com/searKing/sole/web/golang/app/configs/values"
	"github.com/sirupsen/logrus"
)

func MiddlewaresRouter(router gin.IRouter) gin.IRouter {
	logger := logrus.WithField("module", "app.middleware")
	router.Use(gin.LoggerWithWriter(logger.Writer(),
		values.HealthMetricsPrometheusPath,
		values.HealthAliveCheckPath,
		values.HealthReadyCheckPath))
	logger.Infof(`middleware log is loaded`)
	router.Use(gin_.RecoveryWithWriter(logger.Writer(), nil))
	logger.Infof(`middleware recovery is loaded`)

	router.Use(gin_.UseHTTPPreflight())
	logger.Infof(`middleware http preflight is loaded`)
	router.Use(provider.GlobalProvider().GetCORS())
	logger.Infof(`middleware cors is loaded`)

	return router
}
