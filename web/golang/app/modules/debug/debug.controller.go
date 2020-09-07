// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debug

import (
	_ "expvar"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/gin-gonic/gin"
	gin2 "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
)

type DebugController struct {
	pathPrefixTrim string
}

func NewDebugController(prefix string) *DebugController {
	return &DebugController{pathPrefixTrim: prefix}
}

func (d *DebugController) PProf() gin.HandlerFunc {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	if d.pathPrefixTrim != "" {
		return gin2.RedirectTrim(http.StatusFound, d.pathPrefixTrim)
	}
	return gin.WrapH(http.DefaultServeMux)
}

func (d *DebugController) ExpVar() gin.HandlerFunc {
	if d.pathPrefixTrim != "" {
		return gin2.RedirectTrim(http.StatusFound, d.pathPrefixTrim)
	}
	return gin.WrapH(http.DefaultServeMux)
}
