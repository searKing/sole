// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package soletemplate

import (
	"context"

	"github.com/google/wire"

	"github.com/searKing/golang/pkg/webserver"

	"github.com/searKing/sole/pkg/web/modules/date"
	"github.com/searKing/sole/pkg/web/modules/debug"
	"github.com/searKing/sole/pkg/web/modules/doc/swagger"
	"github.com/searKing/sole/pkg/web/modules/index"
	"github.com/searKing/sole/pkg/web/modules/webapp"
	"github.com/searKing/sole/soletemplate/web/modules/soletemplate"
)

// WebHandler is a Wire provider set that includes all Services interface
// implementations.
// interface layer
var WebHandler = wire.NewSet(
	NewWebHandlers,             // interface layer-web
	debug.NewController,        // interface layer-debug
	date.NewController,         // interface layer-date
	index.NewController,        // interface layer-index
	swagger.NewController,      // interface layer-swagger
	webapp.NewController,       // interface layer-webapp
	soletemplate.NewController, // interface layer-example
)

func NewWebHandlers(ws *webserver.WebServer,
	c1 *debug.Controller,
	c2 *date.Controller,
	c3 *index.Controller,
	c4 *swagger.Controller,
	c5 *webapp.Controller,
	c6 *soletemplate.Controller) []webserver.WebHandler {
	ws.AddPostStartHookOrDie("web_handler", func(ctx context.Context) error {
		if c6 != nil {
			c6.PreferRegisterHTTPFromEndpoint = ws.PreferRegisterHTTPFromEndpoint
		}
		ws.InstallWebHandlers(c1, c2, c3, c4, c5, c6)
		return nil
	})
	return []webserver.WebHandler{c1, c2, c3, c4, c5}
}
