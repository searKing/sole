// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"log/slog"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"

	slog_ "github.com/searKing/golang/go/log/slog"
	"github.com/searKing/sole/pkg/domain/logging"

	v1 "github.com/searKing/sole/api/protobuf-spec/soletemplate/v1"

	"github.com/searKing/sole/soletemplate/cmd/soletemplate/app/serve/config"
)

// NewConfig returns a new server configuration from viper.
func NewConfig(ctx context.Context, v *viper.Viper) (c *v1.Configuration, err error) {
	spanName := "NewConfig"
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
	return config.NewViperConfig(v).Complete().New()
}
