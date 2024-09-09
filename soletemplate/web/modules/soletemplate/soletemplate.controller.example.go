// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package soletemplate

import (
	"context"
	"log/slog"

	slog_ "github.com/searKing/golang/go/log/slog"
	v1 "github.com/searKing/sole/api/protobuf-spec/soletemplate/v1"
	"github.com/searKing/sole/pkg/domain/errorgoogle"
	"github.com/searKing/sole/pkg/domain/logging"
	"github.com/searKing/sole/pkg/domain/truncate"
	"github.com/searKing/sole/soletemplate/pkg/domain/templateexample"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Example 样例接口
func (c *Controller) Example(ctx context.Context, req *v1.ExampleRequest) (resp *v1.ExampleResponse, err error) {
	spanName := "soletemplate.Example"
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	defer span.End()
	span.SetAttributes(attribute.String("request_id", req.GetRequestId()))
	span.SetAttributes(attribute.String("app_id", req.GetAppId()))
	logger := slog.With(logging.SpanAttrs(span)...).
		With("request_id", req.GetRequestId()).
		With("app_id", req.GetAppId())
	logger = logger.With("cmd", "Example")
	logger.With("req", truncate.DefaultTruncate(req)).Info("received a request")
	defer func() { err = ApiError(err) }()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cmdReq := &templateexample.ExampleRequest{
		Id:      req.GetRequestId(),
		Message: req.GetMessage(),
	}

	cmdResp, err := c.app.Commands.TemplateExample.Handle(ctx, cmdReq)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "TemplateExample failed")
		logger.With(slog_.Error(err)).Error("failed to run command TemplateExample")
		return nil, errorgoogle.Errore(v1.SoleTemplateErrorEnum_INTERNAL, ApiError(err))
	}

	resp = &v1.ExampleResponse{Message: cmdResp.Message}
	logger.With("resp", truncate.DefaultTruncate(resp)).Info("send a response after running commands")
	return resp, nil
}
