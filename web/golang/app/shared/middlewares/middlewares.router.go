// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/searKing/sole/internal/pkg/provider"
)

func MiddlewaresRouter(router gin.IRouter) gin.IRouter {
	logger := logrus.WithField("module", "app.middleware")
	router.Use(provider.GlobalProvider().GetCORS())
	logger.Infof(`middleware cors is loaded`)

	return router
}
