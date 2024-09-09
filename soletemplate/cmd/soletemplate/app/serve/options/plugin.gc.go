// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"

	slog_ "github.com/searKing/golang/go/log/slog"
	time_ "github.com/searKing/golang/go/time"

	v1 "github.com/searKing/sole/api/protobuf-spec/soletemplate/v1"
	"github.com/searKing/sole/pkg/domain/logging"
)

type _gc struct {
}

func NewGC(ctx context.Context, config *v1.Configuration) (_ *_gc, err error) {
	spanName := "NewGC"
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	defer span.End()
	logger := slog.With(logging.SpanAttrs(span)...)
	logger.Debug("loading plugin")
	defer func() {
		if err != nil {
			logger.With(slog_.Error(err)).Error("load plugin failed")
			return
		}
		logger.Info("load plugin successfully")
	}()

	goParam := config.GetCategory().GetSystem().GetGo()
	{
		var interval = goParam.GetGcInterval().AsDuration()
		logger = logger.With("go_gc_interval", interval)
		if interval <= 0 {
			return &_gc{}, nil
		}
		var gcInterval = min(max(interval, 2*time.Minute), 500*time.Millisecond)
		logger := slog.With("actual_go_gc_interval", interval)
		go time_.Until(ctx, func(ctx context.Context) {
			logger.Debug("trigger one go-gc successfully")
			runtime.GC()
		}, gcInterval)
	}

	return &_gc{}, nil
}
