// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	_ "embed"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	gin_ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
	"github.com/searKing/golang/third_party/github.com/gin-gonic/gin/render"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
)

// Controller ...
type Controller struct{}

// NewController ...
func NewController() *Controller {
	return &Controller{}
}

// SetRoutes registers this handler's routes.
func (c *Controller) SetRoutes(ginRouter gin.IRouter, grpcRouter *grpc.Gateway) {
	slog.Info("installing router")

	// see https://golang.org/pkg/net/http/#FileServer
	// redirect .../index.html to .../ for file
	// As a special case, the returned file server redirects any request ending in "/index.html" to the same path, without the final "index.html".
	ginRouter.GET("/index", gin_.Redirect(http.StatusFound, "/"))
	ginRouter.GET("/index.html", gin_.Redirect(http.StatusFound, "/"))

	ginRouter.GET("/", c.Index())
}

//go:embed index.view.tmpl
var indexView string

func (c *Controller) Index() gin.HandlerFunc {
	r := &render.TemplateHTML{
		Name:  "index",
		Texts: []string{indexView},
		Data:  GetTemplateInfo("/webapp", ""),
	}
	return func(c *gin.Context) {
		c.Render(http.StatusOK, r)
	}
}
