// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package verflag defines utility functions to handle command line flags
// related to version of Kubernetes.
package verflag

import (
	"flag"
	"fmt"
	"os"

	"github.com/searKing/sole/pkg/appinfo"
)

const versionFlagName = "version"

var (
	versionFlag = flag.Bool(versionFlagName, false, "Print version information and quit")
)

// AddFlags registers this package's flags on arbitrary FlagSets, such that they point to the
// same value as the global flags.
func AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(versionFlag, versionFlagName, false, "Print version information and quit")
}

// PrintAndExitIfRequested will check if the -version flag was passed
// and, if so, print the version and exit.
func PrintAndExitIfRequested() {
	if *versionFlag {
		fmt.Printf("%s %s\n", appinfo.ServiceName, appinfo.GetVersion().String())
		os.Exit(0)
	}
}
