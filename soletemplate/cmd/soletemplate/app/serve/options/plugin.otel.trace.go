// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-logr/stdr"
	slog_ "github.com/searKing/golang/go/log/slog"
	"github.com/searKing/golang/pkg/webserver"
	"go.opentelemetry.io/otel"

	trace_ "github.com/searKing/golang/pkg/instrumentation/otel/trace"
	v1 "github.com/searKing/sole/api/protobuf-spec/soletemplate/v1"
	"github.com/searKing/sole/pkg/domain/logging"

	_ "github.com/searKing/golang/pkg/instrumentation/otel/trace/otlptrace/otlptracegrpc" // for otel-grpc
	_ "github.com/searKing/golang/pkg/instrumentation/otel/trace/otlptrace/otlptracehttp" // for otel-http
	_ "github.com/searKing/golang/pkg/instrumentation/otel/trace/stdouttrace"             // for stdout
)

type _otelTrace struct{}

func NewOtelTrace(ctx context.Context, ws *webserver.WebServer, cfg *v1.Configuration) (_ *_otelTrace, err error) {
	spanName := "NewOtelTrace"
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	defer span.End()
	logger := slog.With(logging.SpanAttrs(span)...)
	defer func() {
		if err != nil {
			logger.With(slog_.Error(err)).Error("load plugin failed")
			return
		}
		logger.Info("load plugin successfully")
	}()

	otel.SetLogger(stdr.New(slog.NewLogLogger(slog.Default().Handler(), slog.LevelDebug)))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) { slog.With(slog_.Error(err)).Error("otel: handler returned an error") }))

	otelpb := cfg.GetOtel()
	tp, err := trace_.NewTracerProvider(context.Background(),
		trace_.WithOptionExporterEndpoints(otelpb.GetTraceExporterEndpoints()...),
		trace_.WithOptionResourceAttrs(toOtelAttrs(otelpb.GetResourceAttrs())...))
	if err != nil {
		logger.With(slog_.Error(err)).Error("install otel trace, exit")
		return nil, err
	}
	otel.SetTracerProvider(tp)

	ws.AddPreShutdownHookOrDie("otel-trace", func() error {
		err := tp.Shutdown(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
		return nil
	})
	return &_otelTrace{}, nil
}
