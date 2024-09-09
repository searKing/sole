// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"github.com/google/wire"

	"github.com/searKing/sole/soletemplate/pkg/application"
	"github.com/searKing/sole/soletemplate/web/app/soletemplate"
)

// WebHandler is a wire provider set that includes all web service repositories interface implementations.
// service interface layer
var WebHandler = wire.NewSet(
	soletemplate.WebHandler, // interface layer
	Application,             // application layer
)

// Application is a wire provider set that includes all basic repositories interface implementations.
var Application = wire.NewSet(
	application.DefaultApplication,
	NewTemplateExampleRepository,
	NewConfig,
	NewOtelMetric,
	NewOtelTrace,
)
