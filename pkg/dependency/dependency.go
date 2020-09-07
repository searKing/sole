// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dependency

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var errNilDependency = errors.New("A dependency was expected to be defined but is nil. Please open an issue with the stack trace.")

func ExpectDependency(logger logrus.FieldLogger, dependencies map[string]interface{}) {
	for name, d := range dependencies {
		if d == nil {
			logger.WithError(errors.WithStack(errNilDependency)).Fatalf("A fatal issue occurred. module[%s] is missing", name)
		}
	}
}
