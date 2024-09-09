// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errorgoogle

import (
	"fmt"

	"golang.org/x/exp/constraints"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errors_ "github.com/searKing/golang/go/errors"
)

func Errorf[E constraints.Integer](code E, format string, a ...any) error {
	// try to convert to gRPC standard error code
	c := codes.Code(code)
	if c == codes.OK {
		return nil
	}

	var reason = c.String()

	s := status.New(c, fmt.Sprintf(format, a...))
	if s.Code() != codes.OK && reason != "" {
		details, err := s.WithDetails(&epb.ErrorInfo{
			Reason: reason,
		})
		if err == nil {
			s = details
		}
	}
	return s.Err()
}

// Errore converts the given error into a gRPC error with the given code if
// err's type is not *status.Status.
// ignore code if err is *status.Status.
func Errore[E constraints.Integer](code E, err error) error {
	if err == nil {
		return nil
	}
	_, ok := status.FromError(err)
	if ok {
		return err
	}
	return errors_.Mark(Errorf(code, err.Error()), err)
}
