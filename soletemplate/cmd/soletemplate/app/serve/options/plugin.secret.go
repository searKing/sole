// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"log/slog"

	v1 "github.com/searKing/sole/api/protobuf-spec/soletemplate/v1"
	"github.com/searKing/sole/pkg/domain/logging"
	"go.opentelemetry.io/otel"
	"gocloud.dev/secrets"
	_ "gocloud.dev/secrets/localsecrets"

	slog_ "github.com/searKing/golang/go/log/slog"
)

func NewSecret(ctx context.Context, config *v1.Configuration) (_ *secrets.Keeper, _ func(), err error) {
	spanName := "NewSecret"
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
	keeperURL := config.GetCategory().GetSystem().GetSecretKeeperUrl()
	// Open a *secrets.Keeper using the keeperURL.
	keeper, err := secrets.OpenKeeper(ctx, keeperURL)
	if err != nil {
		return nil, nil, err
	}

	return keeper, func() {
		err := keeper.Close()
		if err != nil {
			slog.With(slog_.Error(err)).With("keeper_url", keeperURL).Error("failed to close secrets keeper")
			return
		}
	}, nil
}
