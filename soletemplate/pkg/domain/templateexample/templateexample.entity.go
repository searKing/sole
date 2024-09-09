// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templateexample

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/searKing/sole/pkg/domain/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var _ Repository = (*TemplateExample)(nil)

// TemplateExample 服务模板
type TemplateExample struct {
	Id string

	CreatedAt time.Time
}

type ExampleRequest struct {
	Id      string
	Message string
}

type ExampleResponse struct {
	Message string
}

var ErrMessageEmpty = errors.New("empty message")

func (e *TemplateExample) Example(ctx context.Context, req *ExampleRequest) (resp *ExampleResponse, err error) {
	ctx, span := otel.Tracer("").Start(ctx, "TemplateExample.Example")
	defer span.End()
	span.SetAttributes(attribute.String("id", req.Id))
	logger := slog.With(logging.SpanAttrs(span)...)

	logger.Info("TemplateExample executed")
	return &ExampleResponse{
		Message: fmt.Sprintf("Get %s at %s", req.Message, time.Now()),
	}, nil
}
