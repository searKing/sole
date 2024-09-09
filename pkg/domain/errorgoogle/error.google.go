// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errorgoogle

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/searKing/sole/api/protobuf-spec/sole/types/v1/errors"
)

func NewErrorStatus(err error) *errors.Error_Status {
	s, ok := status.FromError(err)
	if !ok {
		if err == nil {
			s = status.New(codes.OK, codes.OK.String())
		} else {
			s = status.New(codes.Unknown, err.Error())
		}
	}
	var es errors.Error_Status
	es.Code = int32(s.Code())
	es.Message = s.Message()
	es.Details = s.Proto().GetDetails()
	return &es
}
