// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debug

import (
	"github.com/gin-gonic/gin"
	"github.com/searKing/sole/web/golang/app/configs/values"
)

func DebugRouter(router gin.IRouter, prefix string) gin.IRouter {
	debug := NewDebugController(prefix)
	router.GET(values.DebugPProf, debug.PProf())
	router.GET(values.DebugExpVar, debug.ExpVar())

	return router
}
