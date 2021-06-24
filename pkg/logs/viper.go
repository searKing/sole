// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logs

import (
	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
	"github.com/searKing/sole/pkg/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// RegisterVipers adds vipers functions to the given scheme.
// Public to allow building arbitrary schemes.
// All generated defaulters are covering - they call all nested defaulters.
func RegisterVipers(scheme *runtime.Scheme) error {
	scheme.AddTypeViperLoadingFunc(&Config{}, func(obj interface{}, v *viper.Viper) error { return obj.(*Config).SetViperConfig(v) })
	return nil
}

// SetViperConfig assigns values from viper for the Config
func (c *Config) SetViperConfig(v *viper.Viper) error {
	if c.GetViper != nil {
		v = c.GetViper()
	}
	if err := viper_.UnmarshalProtoMessageByJsonpb(v, &c.Proto); err != nil {
		logrus.WithError(err).Errorf("load logs config from viper")
		return err
	}
	return nil
}
