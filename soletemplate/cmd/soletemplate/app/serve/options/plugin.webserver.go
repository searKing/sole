// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"go.opentelemetry.io/otel"

	slices_ "github.com/searKing/golang/go/exp/slices"
	slog_ "github.com/searKing/golang/go/log/slog"
	"github.com/searKing/golang/go/version"
	"github.com/searKing/golang/pkg/webserver"

	v1 "github.com/searKing/sole/api/protobuf-spec/soletemplate/v1"
	"github.com/searKing/sole/pkg/domain/logging"
	webserver_ "github.com/searKing/sole/pkg/domain/webserver"
	"github.com/searKing/sole/soletemplate/cmd/soletemplate/app/serve/config"
)

func NewWebServer(ctx context.Context, cfg *v1.Configuration) (ws *webserver.WebServer, err error) {
	spanName := "NewWebServer"
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
	web := cfg.GetWeb()

	var fc webserver.FactoryConfig
	fc.SetDefaults()
	fc.Name = version.ServiceName
	fc.BindAddress = net.JoinHostPort(web.GetBindAddr().GetHost(), fmt.Sprintf("%d", web.GetBindAddr().GetPort()))
	fc.PreferRegisterHTTPFromEndpoint = web.GetPreferRegisterHttpFromEndpoint()
	fc.ForceDisableTls = slices_.Or(config.ForceDisableTls, web.GetForceDisableTls())
	fc.ShutdownDelayDuration = web.GetShutdownDelayDuration().AsDuration()

	mw := web.GetMiddlewares()
	fc.MaxConcurrencyUnary = int(mw.GetMaxConcurrencyUnary())
	fc.MaxConcurrencyStream = int(mw.GetMaxConcurrencyStream())
	fc.BurstLimitTimeoutUnary = mw.GetBurstLimitTimeoutUnary().AsDuration()
	fc.BurstLimitTimeoutStream = mw.GetBurstLimitTimeoutStream().AsDuration()
	fc.HandledTimeoutUnary = mw.GetHandledTimeoutUnary().AsDuration()
	fc.HandledTimeoutStream = mw.GetHandledTimeoutStream().AsDuration()
	fc.MaxReceiveMessageSizeInBytes = int(mw.GetMaxReceiveMessageSizeInBytes())
	fc.MaxSendMessageSizeInBytes = int(mw.GetMaxSendMessageSizeInBytes())
	fc.StatsHandling = mw.GetStatsHandling()
	fc.FillRequestId = mw.GetFillRequestId()
	fc.OtelHandling = mw.GetOtelHandling()

	fc.GatewayOptions = append(fc.GatewayOptions, webserver_.GatewayOptions()...)
	return webserver.NewWebServer(fc)
}
