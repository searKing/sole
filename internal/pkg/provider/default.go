// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"github.com/google/uuid"

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

	proto.Service = &viper.Service{}
	proto.GetService().Name = version.ServiceName
	proto.GetService().DisplayName = version.ServiceName
	proto.GetService().Description = version.ServiceDescription
	proto.GetService().Id = proto.GetService().GetName() + "-" + uuid.New().String()

	return proto
}
