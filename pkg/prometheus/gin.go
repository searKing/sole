// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prometheus

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	strings_ "github.com/searKing/golang/go/strings"
	"github.com/sirupsen/logrus"
)

const (
	namespace = "sole" // For prefixing prometheus metrics
)

var (
	openConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "http_connections_open",
		Help:      "Current number of http open connections.",
	})

	httpRequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_requests_total",
			Help:      "Description",
		},
		[]string{"endpoint"},
	)
	httpResponseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_response_time_seconds",
			Help:      "Description",
		},
		[]string{"endpoint"},
	)
)

func GinHttpMetric(notMetric ...string) gin.HandlerFunc {
	logrus.Info("gin http metric registered")
	return func(c *gin.Context) {
		if strings_.SliceContains(notMetric, c.Request.URL.Path) {
			c.Next()
			return
		}

		start := time.Now()
		openConnections.Inc()
		c.Next()
		//Request counter metric
		httpRequestCounter.WithLabelValues(c.Request.URL.Path).Inc()
		//Response time metric
		httpResponseTime.WithLabelValues(c.Request.URL.Path).Observe(time.Since(start).Seconds())
		openConnections.Dec()

	}
}
