// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"github.com/gin-gonic/gin"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway/grpc"
	"github.com/searKing/sole/web/golang"
)

func Setup() (*gin.Engine, *grpc.Gateway) {
	ginBackend := setupHTTP()
	grpcBackend := setupGRPC()
	grpcBackend.Handler = ginBackend

	golang.NewHandler().SetRoutes(ginBackend, grpcBackend)
	return ginBackend, grpcBackend
}
