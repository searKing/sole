// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpReqs = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "date_http_request_total",
			Help: "How many Date HTTP requests processed.",
		},
		[]string{},
	)
)

var (
	httpReqCostInMilliSeconds = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "date_cost_duration_milliseconds",
			Help: "The duration of the date cost in milliseconds.",
		})
)
