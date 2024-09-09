// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package swagger

import (
	_ "embed"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/searKing/sole/api/openapi-spec/sole"

	"github.com/searKing/golang/go/version"
	"github.com/searKing/golang/third_party/github.com/gin-gonic/gin/render"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
)

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

// SetRoutes registers this handler's routes.
func (c *Controller) SetRoutes(ginRouter gin.IRouter, grpcRouter *grpc.Gateway) {
	slog.Info("installing router")
	ginRouter.GET("/doc/swagger/swagger.json", c.Content("swagger.json", strings.NewReader(v1.SwaggerUIJson)))
	ginRouter.GET("/doc/swagger/swagger.yaml", c.Content("swagger.yaml", strings.NewReader(v1.SwaggerUIYaml)))

	for _, path_ := range []string{
		"/doc/swagger/swagger-ui",
		"/doc/swagger/swagger-ui/index.html",
	} {
		ginRouter.GET(path_, c.UI())
	}
}

func (c *Controller) Content(name string, content io.ReadSeeker) gin.HandlerFunc {
	return func(ctx *gin.Context) { http.ServeContent(ctx.Writer, ctx.Request, name, time.Now(), content) }
}

//go:embed swagger.view.index.tmpl
var swaggerUIView string

func (c *Controller) UI() gin.HandlerFunc {
	r := &render.TemplateHTML{
		Name:  "swagger-ui",
		Texts: []string{swaggerUIView},
		Data: struct {
			Name           string
			BaseUrl        string
			SwaggerJsonUrl string
		}{
			Name:           version.ServiceName,
			SwaggerJsonUrl: "/doc/swagger/swagger.json",
		},
	}
	return func(c *gin.Context) {
		c.Render(http.StatusOK, r)
	}
}
