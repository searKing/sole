// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webapp

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"

	http_ "github.com/searKing/sole/pkg/domain/http"
)

var RelativeFileStoragePathPrefixes = []string{"./web/webapp", "../pkg/web/webapp"}

// Controller ...
type Controller struct{}

// NewController ...
func NewController() *Controller {
	return &Controller{}
}

// SetRoutes registers this handler's routes.
func (c *Controller) SetRoutes(ginRouter gin.IRouter, grpcRouter *grpc.Gateway) {
	slog.Info("installing router")

	ginRouter.StaticFS("/webapp", http_.Dirs(RelativeFileStoragePathPrefixes))
}
