// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	providerReloads = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "provider_reload_total",
			Help: "How many provider reloaded.",
		}, []string{"proto"})
)
