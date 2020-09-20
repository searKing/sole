// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viper

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	filepath_ "github.com/searKing/golang/go/path/filepath"
	"github.com/searKing/sole/api/protobuf-spec/v1/viper"
)

var (
	Version         = "dev-master"
	BuildTime       = "undefined"
	GitHash         = "undefined"
	ForceDisableTls bool
)

const (
	ServiceName        = "sole"
	ServiceDescription = "sole is a cloud native high throughput service manager server, " +
		"allowing you to manage all services."
)

func NewDefaultViperProto() *viper.ViperProto {
	proto := &viper.ViperProto{}

	proto.AppInfo = &viper.AppInfo{}
	proto.GetAppInfo().BuildVersion = Version
	proto.GetAppInfo().BuildTime = BuildTime
	proto.GetAppInfo().BuildHash = GitHash

	proto.Service = &viper.Service{}
	proto.GetService().Name = ServiceName
	proto.GetService().DisplayName = ServiceName
	proto.GetService().Description = ServiceDescription
	proto.GetService().Id = uuid.New().String()

	proto.Log = &viper.Log{}
	proto.GetLog().Level = viper.Log_info
	proto.GetLog().Format = viper.Log_text
	proto.GetLog().Path = "./log/" + ServiceName
	proto.GetLog().RotationDuration = ptypes.DurationProto(24 * time.Hour)
	proto.GetLog().RotationMaxCount = 0
	proto.GetLog().RotationMaxAge = ptypes.DurationProto(7 * 24 * time.Hour)

	proto.Web = &viper.Web{}
	proto.GetWeb().ForceDisableTls = ForceDisableTls
	proto.GetWeb().BindAddr = &viper.Web_Net{}
	if proto.GetWeb().ForceDisableTls {
		proto.GetWeb().GetBindAddr().Port = 80
	} else {
		proto.GetWeb().GetBindAddr().Port = 443
	}

	proto.Database = &viper.Database{}
	proto.GetDatabase().Dsn = "memory"

	return proto
}

// DefaultConfigPath returns config file's default path
func DefaultConfigPath() string {
	return filepath_.Pathify(fmt.Sprintf("$HOME/.%s.yaml", ServiceName))
}
