// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag

import (
	"flag"

	"github.com/sirupsen/logrus"
)

// InitFlags normalizes, parses, then logs the command line flags
func InitFlags() {
	flag.Parse()
	PrintFlags(flag.CommandLine)
}

// PrintFlags logs the flags in the flagset
func PrintFlags(flags *flag.FlagSet) {
	flags.VisitAll(func(flag *flag.Flag) {
		logrus.Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}
