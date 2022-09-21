// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"runtime"
	"time"

	math_ "github.com/searKing/golang/go/exp/math"
	time_ "github.com/searKing/golang/go/time"
	configpb "github.com/searKing/sole/api/protobuf-spec/v1/config"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type _gc struct {
}

func NewGC(ctx context.Context, config *configpb.Configuration) (_ *_gc, err error) {
	spanName := "NewGC"
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	defer span.End()
	logger := logrus.WithField("trace_id", span.SpanContext().TraceID()).
		WithField("span_id", span.SpanContext().SpanID())
	defer func() {
		if err != nil {
			logger.WithError(err).Error("load plugin failed")
			return
		}
		logger.Info("load plugin successfully")
	}()

	goParam := config.GetGo()
	{
		var interval = goParam.GetGcInterval().AsDuration()
		logger = logger.WithField("go_gc_interval", interval)
		if interval <= 0 {
			return &_gc{}, nil
		}
		var gcInterval = math_.Min(math_.Max(interval, 2*time.Minute), 500*time.Millisecond)
		logger := logrus.WithField("actual_go_gc_interval", interval)
		go time_.Until(ctx, func(ctx context.Context) {
			logger.Debugf("trigger one go-gc successfully")
			runtime.GC()
		}, gcInterval)
	}

	return &_gc{}, nil
}
