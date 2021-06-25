// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logs

import (
	"os"
	"path/filepath"
	"time"

	"github.com/searKing/sole/pkg/runtime"
	"google.golang.org/protobuf/types/known/durationpb"
)

// RegisterDefaults adds defaulters functions to the given scheme.
// Public to allow building arbitrary schemes.
// All generated defaulters are covering - they call all nested defaulters.
func RegisterDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&Config{}, func(obj interface{}) { obj.(*Config).SetDefaultsConfig() })
	return nil
}

// SetDefaultsConfig assigns default values for the Config
func (c *Config) SetDefaultsConfig() *Config {
	c.Proto = Log{
		Level:              Log_info,
		Format:             Log_text,
		Path:               "./log/" + filepath.Base(os.Args[0]),
		RotationDuration:   durationpb.New(24 * time.Hour),
		RotationMaxCount:   0,
		RotationMaxAge:     durationpb.New(7 * 24 * time.Hour),
		ReportCaller:       false,
		MuteDirectlyOutput: true,
	}
	return c
}
