// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golang

import (
	"context"

	"github.com/google/wire"
	"github.com/searKing/golang/pkg/webserver"
	"github.com/searKing/sole/web/golang/app/modules/date"
	"github.com/searKing/sole/web/golang/app/modules/debug"
	"github.com/searKing/sole/web/golang/app/modules/doc/swagger"
	"github.com/searKing/sole/web/golang/app/modules/index"
	"github.com/searKing/sole/web/golang/app/modules/webapp"
)

// WebHandler is a Wire provider set that includes all Services interface
// implementations.
var WebHandler = wire.NewSet(
	NewWebHandlers,      // 接口层
	debug.NewController, // 接口层-调试
	date.NewController,  // 接口层-日期
	index.NewController,
	webapp.NewController,
	swagger.NewController,
)

func NewWebHandlers(ws *webserver.WebServer, c1 *debug.Controller, c2 *date.Controller,
	c3 *index.Controller, c4 *webapp.Controller, c5 *swagger.Controller) []webserver.WebHandler {
	ws.AddPostStartHookOrDie("web_handler", func(ctx context.Context) error {
		ws.InstallWebHandlers(c1, c2, c3, c4, c5)
		return nil
	})
	return []webserver.WebHandler{c1, c2, c3, c4, c5}
}
