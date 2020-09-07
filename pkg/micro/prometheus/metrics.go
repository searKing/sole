// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// Metrics prototypes
type Metrics struct {
	OpenConnections prometheus.Gauge

	HTTPRequestCounter *prometheus.CounterVec
	HTTPResponseTime   *prometheus.HistogramVec
}

// Method for creation new custom Prometheus  metrics
func NewMetrics(name, version, hash, buildTime string) *Metrics {
	pm := &Metrics{
		OpenConnections: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_http_connections_open", name),
			Help: "Current number of http open connections.",
		}),
		HTTPRequestCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%s_http_requests_total", name),
				Help: "Description",
				ConstLabels: map[string]string{
					"version":   version,
					"hash":      hash,
					"buildTime": buildTime,
				},
			},
			[]string{"endpoint"},
		),
		HTTPResponseTime: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: fmt.Sprintf("%s_http_response_time_seconds", name),
				Help: "Description",
				ConstLabels: map[string]string{
					"version":   version,
					"hash":      hash,
					"buildTime": buildTime,
				},
			},
			[]string{"endpoint"},
		),
	}
	prometheus.MustRegister(pm.OpenConnections)
	prometheus.MustRegister(pm.HTTPRequestCounter)
	prometheus.MustRegister(pm.HTTPResponseTime)
	return pm
}
