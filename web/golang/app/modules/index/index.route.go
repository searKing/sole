// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gin_ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
	"github.com/searKing/sole/web/golang/app/configs/values"
)

func Router(router gin.IRouter) gin.IRouter {
	// see https://golang.org/pkg/net/http/#FileServer
	// redirect .../index.html to .../ for file
	// As a special case, the returned file server redirects any request ending in "/index.html" to the same path, without the final "index.html".
	router.GET(values.Index, gin_.Redirect(http.StatusFound, values.IndexAsBase))
	router.GET(values.IndexAsHtml, gin_.Redirect(http.StatusFound, values.IndexAsBase))

	router.GET(values.IndexAsBase, Controller())
	return router
}
