// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package soletemplate

import (
	"context"

	grpc_ "github.com/searKing/golang/third_party/google.golang.org/grpc"

	v1 "github.com/searKing/sole/api/protobuf-spec/soletemplate/v1"
)

type HttpController struct {
	*Controller

	decorators grpc_.UnaryHandlerDecorators
}

func (c *HttpController) Health(ctx context.Context, req *v1.HealthRequest) (resp *v1.HealthResponse, err error) {
	return grpc_.WithUnaryHandlerDecorators(c.Controller.Health, c.decorators...)(ctx, req)
}

func (c *HttpController) Encrypt(ctx context.Context, req *v1.EncryptRequest) (*v1.EncryptResponse, error) {
	return grpc_.WithUnaryHandlerDecorators(c.Controller.Encrypt, c.decorators...)(ctx, req)
}

func (c *HttpController) Example(ctx context.Context, req *v1.ExampleRequest) (*v1.ExampleResponse, error) {
	return grpc_.WithUnaryHandlerDecorators(c.Controller.Example, c.decorators...)(ctx, req)
}
