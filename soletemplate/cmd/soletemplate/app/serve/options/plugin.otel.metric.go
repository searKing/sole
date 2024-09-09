// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-logr/stdr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	maps_ "github.com/searKing/golang/go/exp/maps"
	slog_ "github.com/searKing/golang/go/log/slog"
	"github.com/searKing/golang/pkg/webserver"

	metric_ "github.com/searKing/golang/pkg/instrumentation/otel/metric"
	v1 "github.com/searKing/sole/api/protobuf-spec/soletemplate/v1"
	"github.com/searKing/sole/pkg/domain/logging"

	_ "github.com/searKing/golang/pkg/instrumentation/otel/metric/otlpmetric/otlpmetricgrpc" // for otlp-grpc
	_ "github.com/searKing/golang/pkg/instrumentation/otel/metric/otlpmetric/otlpmetrichttp" // for otlp-http
	_ "github.com/searKing/golang/pkg/instrumentation/otel/metric/prometheusmetric"          // for prometheus
	_ "github.com/searKing/golang/pkg/instrumentation/otel/metric/stdoutmetric"              // for stdout
)

type _otelMetric struct{}

func NewOtelMetric(ctx context.Context, ws *webserver.WebServer, cfg *v1.Configuration) (_ *_otelMetric, err error) {
	spanName := "NewOtelMetric"
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
	mp, err := metric_.NewMeterProvider(context.Background(),
		metric_.WithOptionExporterEndpoints(otelpb.GetMetricExporterEndpoints()...),
		metric_.WithOptionResourceAttrs(toOtelAttrs(otelpb.GetResourceAttrs())...))
	if err != nil {
		logger.With(slog_.Error(err)).Error("install otel metric, exit")
		return nil, err
	}
	otel.SetMeterProvider(mp)

	ws.AddPreShutdownHookOrDie("otel-metric", func() error {
		err := mp.Shutdown(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
		return nil
	})
	return &_otelMetric{}, nil
}

func toOtelAttrs(attrs map[string]string) []attribute.KeyValue {
	return maps_.SliceFunc(attrs, func(k, v string) attribute.KeyValue {
		return attribute.String(k, v)
	})
}
