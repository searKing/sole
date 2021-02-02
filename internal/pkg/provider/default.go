// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	filepath_ "github.com/searKing/golang/go/path/filepath"

	"github.com/searKing/sole/api/protobuf-spec/v1/viper"
	"github.com/searKing/sole/internal/pkg/version"
)

var ForceDisableTls bool

func NewDefaultViperProto() *viper.ViperProto {
	proto := &viper.ViperProto{}

	proto.AppInfo = &viper.AppInfo{}
	proto.GetAppInfo().BuildVersion = version.Version
	proto.GetAppInfo().BuildTime = version.BuildTime
	proto.GetAppInfo().BuildHash = version.GitHash
	proto.GetAppInfo().GoVersion = version.GoVersion
	proto.GetAppInfo().Compiler = version.Compiler
	proto.GetAppInfo().Platform = version.Platform

	proto.Service = &viper.Service{}
	proto.GetService().Name = version.ServiceName
	proto.GetService().DisplayName = version.ServiceName
	proto.GetService().Description = version.ServiceDescription
	proto.GetService().Id = proto.GetService().GetName() + "-" + uuid.New().String()

	proto.Log = &viper.Log{}
	proto.GetLog().Level = viper.Log_info
	proto.GetLog().Format = viper.Log_text
	proto.GetLog().Path = "./log/" + version.ServiceName
	proto.GetLog().RotationDuration = ptypes.DurationProto(24 * time.Hour)
	proto.GetLog().RotationMaxCount = 0
	proto.GetLog().RotationMaxAge = ptypes.DurationProto(7 * 24 * time.Hour)
	proto.GetLog().ReportCaller = true

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
	// 	return filepath_.Pathify(fmt.Sprintf("$HOME/.%s.yaml", version.ServiceName))
	return filepath_.Pathify(fmt.Sprintf("./conf/%s.yaml", version.ServiceName))
}
