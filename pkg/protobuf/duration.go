// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protobuf

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/durationpb"
)

func DurationOrDefault(timeout *durationpb.Duration, def time.Duration, msg string) time.Duration {
	if timeout == nil {
		return def
	}
	d, err := ptypes.Duration(timeout)
	if err != nil {
		logrus.WithField("timeout", timeout).
			WithError(err).
			Warnf("malformed %s, use %s instead", msg, def)
		return def
	}
	return d
}
