// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package values

import "github.com/searKing/sole/api/protobuf-spec/v1/health"

var (
	HealthMetricsPrometheusPath = health.Pattern_HealthService_MetricsPrometheus_0.String() // "/health/metrics/prometheus"
	HealthAliveCheckPath        = health.Pattern_HealthService_Alive_0.String()             // "/health/alive"
	HealthReadyCheckPath        = health.Pattern_HealthService_Ready_0.String()             // "/health/ready"
	HealthVersionPath           = health.Pattern_HealthService_Version_0.String()           // "/version"
)
