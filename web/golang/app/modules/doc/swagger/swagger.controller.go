// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package swagger

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/searKing/golang/third_party/github.com/gin-gonic/gin/render"
)

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Json() gin.HandlerFunc {
	const serviceSwaggerJson = "web/webapp/static/swagger/swagger.json"
	return func(ctx *gin.Context) {
		ctx.File(serviceSwaggerJson)
	}
}

func (c *Controller) Yaml() gin.HandlerFunc {
	const serviceSwaggerYaml = "web/webapp/static/swagger/swagger.yaml"
	return func(ctx *gin.Context) {
		ctx.File(serviceSwaggerYaml)
	}
}

func (c *Controller) UI() gin.HandlerFunc {
	const swaggerIndexTmplName = "web/webapp/WEB-INF/views/swagger-ui/index.tmpl"
	r := &render.TemplateHTML{
		Name:  "index",
		Files: []string{swaggerIndexTmplName},
		Data:  GetIndexTemplateInfo(""),
	}
	return func(c *gin.Context) {
		c.Render(http.StatusOK, r)
	}
}
