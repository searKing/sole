// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errorgoogle

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/otel"

	slog_ "github.com/searKing/golang/go/log/slog"
)

// HTTPError translate error into error schema for Google's JSON HTTP APIs.
// https://cloud.google.com/apis/design/errors
func HTTPError(ctx context.Context, mux *runtime.ServeMux,
	marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	spanName := "HTTPError"
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	defer span.End()
	logger := slog.With("trace_id", span.SpanContext().TraceID()).
		With("span_id", span.SpanContext().SpanID())

	// This message defines the error schema for Google's JSON HTTP APIs.
	body := NewErrorStatus(err)
	if err != nil {
		logger.With(slog_.Error(err)).
			With("method", r.Method).
			With("url", r.RequestURI).
			With("error_code", body.GetCode()).
			With("error_message", body.GetMessage()).
			Error("handle http request failed")
	}
	// 200 OK forever
	runtime.ForwardResponseMessage(ctx, mux, marshaler, w, r, body)
}
