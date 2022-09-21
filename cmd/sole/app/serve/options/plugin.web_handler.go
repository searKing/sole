// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"github.com/google/wire"
	"github.com/searKing/sole/web/golang"
)

// WebHandler is a Wire provider set that includes all Services interface
// implementations.
// 人脸驱动服务接口层
var WebHandler = wire.NewSet(
	golang.WebHandler, // 接口层
	Application,       // 应用层
)

// Application 人脸驱动服务应用层
var Application = wire.NewSet(
	NewConfig, // 配置文件
)
