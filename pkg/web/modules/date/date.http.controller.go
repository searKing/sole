// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"context"

	grpc_ "github.com/searKing/golang/third_party/google.golang.org/grpc"

	v1 "github.com/searKing/sole/api/protobuf-spec/sole/date/v1"
)

type HttpController struct {
	*Controller
	
	decorators grpc_.UnaryHandlerDecorators
}

// Now Date Query
func (c *HttpController) Now(ctx context.Context, req *v1.DateRequest) (resp *v1.DateResponse, err error) {
	return grpc_.WithUnaryHandlerDecorators(c.Controller.Now, c.decorators...)(ctx, req)
}

// Error Date Query, only return error, for test only
func (c *HttpController) Error(ctx context.Context, req *v1.DateRequest) (resp *v1.DateResponse, err error) {
	return grpc_.WithUnaryHandlerDecorators(c.Controller.Error, c.decorators...)(ctx, req)
}
