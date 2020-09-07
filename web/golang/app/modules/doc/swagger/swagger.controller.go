// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package swagger

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/searKing/golang/third_party/github.com/gin-gonic/gin/render"
	"github.com/searKing/sole/web/golang/app/modules/webapp"
)

type SwaggerController struct{}

func NewSwaggerController() *SwaggerController {
	return &SwaggerController{}
}

func (c *SwaggerController) Json() gin.HandlerFunc {
	const serviceSwaggerJson = "web/webapp/static/swagger/swagger.json"
	return func(ctx *gin.Context) {
		ctx.File(serviceSwaggerJson)
	}
}

func (c *SwaggerController) Yaml() gin.HandlerFunc {
	const serviceSwaggerYaml = "web/webapp/static/swagger/swagger.yaml"
	return func(ctx *gin.Context) {
		ctx.File(serviceSwaggerYaml)
	}
}

func (c *SwaggerController) UI() gin.HandlerFunc {
	const swaggerIndexTmplName = "web/webapp/WEB-INF/views/swagger-ui/index.tmpl"
	const swaggerIndexHtmlName = "web/webapp/app/modules/swagger-ui/index.html"

	return func(c *gin.Context) {
		c.Render(http.StatusOK, render.TemplateHTML{
			Name:  "index",
			Files: []string{swaggerIndexTmplName},
			Data:  GetIndexTemplateInfo(webapp.ResolveWeb(swaggerIndexHtmlName)),
		})
	}
}
