// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package swagger

import (
	"github.com/gin-gonic/gin"
	"github.com/searKing/sole/web/golang/app/configs/values"
)

func SwaggerRouter(router gin.IRouter) gin.IRouter {
	s := NewSwaggerController()
	router.GET(values.SwaggerJson, s.Json())
	router.GET(values.SwaggerYaml, s.Yaml())

	for _, path := range values.SwaggerUis {
		router.GET(path, s.UI())
	}
	return router
}
