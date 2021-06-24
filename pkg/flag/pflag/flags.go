// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Code borrowed from https://github.com/kubernetes/kubernetes

package pflag

import (
	goflag "flag"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// WordSepNormalizeFunc changes all flags that contain "_" separators
func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}

// WarnWordSepNormalizeFunc changes and warns for flags that contain "_" separators
func WarnWordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		nname := strings.Replace(name, "_", "-", -1)
		logrus.Warningf("%s is DEPRECATED and will be removed in a future version. Use %s instead.", name, nname)

		return pflag.NormalizedName(nname)
	}
	return pflag.NormalizedName(name)
}

// InitFlags normalizes, parses, then logs the command line flags
func InitFlags() {
	pflag.CommandLine.SetNormalizeFunc(WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	pflag.Parse()
	pflag.VisitAll(func(flag *pflag.Flag) {
		logrus.Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}

// PrintFlags logs the flags in the flagset
func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		logrus.Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}
