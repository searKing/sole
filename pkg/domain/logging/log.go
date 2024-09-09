// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logging

import (
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

func SpanAttrs(span trace.Span) []any {
	var attrs []any
	if span.SpanContext().HasTraceID() {
		attrs = append(attrs, slog.Any("trace_id", span.SpanContext().TraceID()))
	}
	if span.SpanContext().HasSpanID() {
		attrs = append(attrs, slog.Any("span_id", span.SpanContext().SpanID()))
	}
	return attrs
}
