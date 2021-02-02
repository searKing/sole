// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viper

import (
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

// Load load config from file and protos into v, and save to a using file
// load sequence: protos..., file, env, replace if member has been set
// that is, later cfg appeared, higher priority cfg has
func Load(v proto.Message, cfgFile string, envPrefix string, protos ...proto.Message) error {
	if err := MergeAll(cfgFile, envPrefix, protos...); err != nil {
		logrus.WithError(err).Fatalf("load config proto from the file")
		return err
	}
	if err := PersistConfig(); err != nil {
		logrus.WithError(err).Warnf("persist config proto from the file, ignore")
	}

	return Unmarshal(v)
}
