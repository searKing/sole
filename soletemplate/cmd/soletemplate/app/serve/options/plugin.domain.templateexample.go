// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel"

	slog_ "github.com/searKing/golang/go/log/slog"

	"github.com/searKing/sole/pkg/domain/logging"
	"github.com/searKing/sole/soletemplate/pkg/domain/templateexample"
)

func NewTemplateExampleRepository(ctx context.Context, validator *validator.Validate) (f templateexample.Factory, err error) {
	spanName := "NewTemplateExampleRepository"
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
	var fc templateexample.FactoryConfig
	fc.SetDefaults()
	fc.Validator = validator
	return templateexample.NewFactory(fc)
}
