// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gin_ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
	"github.com/searKing/golang/third_party/github.com/gin-gonic/gin/render"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
	"github.com/searKing/sole/web/golang/app/configs/values"
	"github.com/sirupsen/logrus"
)

// Controller ...
type Controller struct{}

// NewController ...
func NewController() *Controller {
	return &Controller{}
}

// SetRoutes registers this handler's routes.
func (c *Controller) SetRoutes(ginRouter gin.IRouter, grpcRouter *grpc.Gateway) {
	logrus.Info("installing router")

	// see https://golang.org/pkg/net/http/#FileServer
	// redirect .../index.html to .../ for file
	// As a special case, the returned file server redirects any request ending in "/index.html" to the same path, without the final "index.html".
	ginRouter.GET(values.Index, gin_.Redirect(http.StatusFound, values.IndexAsBase))
	ginRouter.GET(values.IndexAsHtml, gin_.Redirect(http.StatusFound, values.IndexAsBase))

	ginRouter.GET(values.IndexAsBase, c.Index())
}

// Index ...
func (c *Controller) Index() gin.HandlerFunc {

	const IndexTmplName = "web/webapp/WEB-INF/views/index/index.tmpl"

	r := &render.TemplateHTML{
		Name:  "index",
		Files: []string{IndexTmplName},
		Data:  GetTemplateInfo(values.WebApp, ""),
	}
	return func(c *gin.Context) {
		c.Render(http.StatusOK, r)
	}
}
