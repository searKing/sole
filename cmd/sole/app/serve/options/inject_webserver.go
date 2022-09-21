// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wireinject
// +build wireinject

package options

import (
	"context"

	"github.com/google/wire"
	"github.com/searKing/golang/pkg/webserver"
	"github.com/spf13/viper"
)

// RunServer 人脸驱动应用
func RunServer(ctx context.Context) (_ *_RunningServer, cleanup func(), err error) {
	// This will be filled in by Wire with providers from the provider sets in
	// wire.Build.
	wire.Build(
		_NewRunningServer,
		NewWebServer,   // HTTP+gRPC服务
		viper.GetViper, // 配置
		NewEnv,         // 环境变量配置
		NewLog,         // 日志
		NewFileCleaner, // 文件清理
		NewGC,          // 文件清理
		WebHandler,     // 人脸驱动接口层
	)
	return nil, nil, nil
}

type _RunningServer struct{}

// NewRunningServer 加载配置、启动服务
func _NewRunningServer(ctx context.Context, _ *_env, _ *_log,
	ws *webserver.WebServer, _ *_fileCleaner, _ []webserver.WebHandler, _ *_gc) (s *_RunningServer, err error) {
	prepared, err := ws.PrepareRun()
	if err != nil {
		return nil, err
	}

	return nil, prepared.Run(ctx)
}
