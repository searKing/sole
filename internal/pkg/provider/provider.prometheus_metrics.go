// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"sync/atomic"

	"github.com/searKing/sole/pkg/micro/prometheus"
)

//go:generate go-atomicvalue -type "prometheusMetricsManager<*github.com/searKing/sole/app/shared/micro/prometheus.MetricsManager>"
type prometheusMetricsManager atomic.Value

func (p *Provider) PrometheusMetricsManager() *prometheus.MetricsManager {
	return p.prometheusMetricsManager.Load()
}

func (p *Provider) updatePrometheusMetricsManager() {
	proto := p.Proto()
	logger := p.Logger().WithField("module", "provider.middleware.prometheus")
	appInfo := proto.GetAppInfo()
	logger.Info("Setting up Prometheus middleware")
	p.prometheusMetricsManager.Store(prometheus.NewMetricsManager(proto.GetService().GetName(),
		appInfo.GetBuildVersion(), appInfo.GetBuildHash(), appInfo.GetBuildTime()))
}
