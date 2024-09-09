// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"

	slog_ "github.com/searKing/golang/go/log/slog"

	v1 "github.com/searKing/sole/api/protobuf-spec/soletemplate/v1"
	"github.com/searKing/sole/pkg/domain/logging"
)

type _env struct{}

func NewEnv(ctx context.Context, config *v1.Configuration) (_ *_env, err error) {
	spanName := "NewEnv"
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
	for k, v := range config.GetCategory().GetDynamicEnvironments() {
		err := os.Setenv(k, v)
		if err != nil {
			return nil, err
		}
	}

	return &_env{}, nil
}
