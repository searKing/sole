// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prometheus

import (
	"net/http"
	"time"
)

type MetricsManager struct {
	prometheusMetrics *Metrics
}

func NewMetricsManager(name, version, hash, buildTime string) *MetricsManager {
	return &MetricsManager{
		prometheusMetrics: NewMetrics(name, version, hash, buildTime),
	}
}

// Main middleware method to collect metrics for Prometheus.
func (pmm *MetricsManager) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	pmm.prometheusMetrics.OpenConnections.Inc()
	next(rw, r)
	//Request counter metric
	pmm.prometheusMetrics.HTTPRequestCounter.WithLabelValues(r.URL.Path).Inc()
	//Response time metric
	pmm.prometheusMetrics.HTTPResponseTime.WithLabelValues(r.URL.Path).Observe(time.Since(start).Seconds())
	go pmm.prometheusMetrics.OpenConnections.Dec()
}
