// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tracing

import (
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/searKing/golang/go/error/exception"
	jeagerConf "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/zipkin"
	"github.com/urfave/negroni"
)

//go:generate stringer -type Type -trimprefix=Type
//go:generate jsonenums -type Type
type Type int

const (
	TypeJeager Type = iota
	TypeZipkin
	TypeButt
)

// Tracer encapsulates tracing abilities.
type Tracer struct {
	Type Type
	jeagerConf.Configuration

	closer io.Closer
}

func New() *Tracer {
	return &Tracer{
		Type:          TypeJeager,
		Configuration: jeagerConf.Configuration{},
	}
}

// Setup sets up the tracer. Currently supports jaeger.
func (t *Tracer) Setup() error {
	if t.Type >= TypeButt {
		return exception.NewIllegalStateException1(fmt.Sprintf("unknown tracer: %s", t.Type))
	}
	var configs []jeagerConf.Option
	if t.Type == TypeZipkin {
		zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
		configs = append(
			configs,
			jeagerConf.Injector(opentracing.HTTPHeaders, zipkinPropagator),
			jeagerConf.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		)
	}

	var err error
	t.closer, err = t.Configuration.InitGlobalTracer(t.ServiceName, configs...)
	return err
}

// Close closes the tracer.
func (t *Tracer) Close() error {
	if t == nil {
		return nil
	}
	if t.closer != nil {
		return t.closer.Close()
	}
	return nil
}

func (t *Tracer) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var span opentracing.Span
	opName := r.URL.Path

	// It's very possible that Hydra is fronted by a proxy which could have initiated a trace.
	// If so, we should attempt to join it.
	remoteContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header),
	)

	if err != nil {
		span = opentracing.StartSpan(opName)
	} else {
		span = opentracing.StartSpan(opName, opentracing.ChildOf(remoteContext))
	}

	defer span.Finish()

	r = r.WithContext(opentracing.ContextWithSpan(r.Context(), span))

	next(rw, r)

	ext.HTTPMethod.Set(span, r.Method)
	if negroniWriter, ok := rw.(negroni.ResponseWriter); ok {
		statusCode := uint16(negroniWriter.Status())
		if statusCode >= http.StatusBadRequest {
			ext.Error.Set(span, true)
		}
		ext.HTTPStatusCode.Set(span, statusCode)
	}
}
