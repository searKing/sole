// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"fmt"
	"net"

	slices_ "github.com/searKing/golang/go/exp/slices"
	"github.com/searKing/golang/go/version"
	"github.com/searKing/golang/pkg/webserver"
	configpb "github.com/searKing/sole/api/protobuf-spec/v1/config"
	"github.com/searKing/sole/cmd/sole/app/serve/config"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

func NewWebServer(ctx context.Context, cfg *configpb.Configuration) (ws *webserver.WebServer, err error) {
	spanName := "NewWebServer"
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	defer span.End()
	logger := logrus.WithField("trace_id", span.SpanContext().TraceID()).
		WithField("span_id", span.SpanContext().SpanID())
	defer func() {
		if err != nil {
			logger.WithError(err).Error("load plugin failed")
			return
		}
		logger.Info("load plugin successfully")
	}()
	web := cfg.GetWeb()

	var fc webserver.FactoryConfig
	fc.SetDefaults()
	fc.Name = version.ServiceName
	fc.BindAddress = net.JoinHostPort(web.GetBindAddr().GetHost(), fmt.Sprintf("%d", web.GetBindAddr().GetPort()))
	fc.ForceDisableTls = slices_.Or(config.ForceDisableTls, web.GetForceDisableTls())
	fc.ShutdownDelayDuration = web.GetShutdownDelayDuration().AsDuration()

	fc.MaxConcurrencyUnary = int(web.GetMaxConcurrencyUnary())
	fc.MaxConcurrencyStream = int(web.GetMaxConcurrencyStream())
	fc.BurstLimitTimeoutUnary = web.GetBurstLimitTimeoutUnary().AsDuration()
	fc.BurstLimitTimeoutStream = web.GetBurstLimitTimeoutStream().AsDuration()
	fc.HandledTimeoutUnary = web.GetHandledTimeoutUnary().AsDuration()
	fc.HandledTimeoutStream = web.GetHandledTimeoutStream().AsDuration()
	fc.MaxReceiveMessageSizeInBytes = int(web.GetMaxReceiveMessageSizeInBytes())
	fc.MaxSendMessageSizeInBytes = int(web.GetMaxSendMessageSizeInBytes())
	return webserver.NewWebServer(fc)
}
