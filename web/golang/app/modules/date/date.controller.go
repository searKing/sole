// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"context"
	"time"

	"github.com/searKing/sole/api/protobuf-spec/v1/date"
)

type Controller struct{}

// 日期查询
func (c *Controller) Now(ctx context.Context, req *date.DateRequest) (resp *date.DateResponse, err error) {
	return &date.DateResponse{
		RequestId: req.GetRequestId(),
		Date:      time.Now().String(),
	}, nil
}
