// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/searKing/golang/third_party/github.com/gin-gonic/gin/render"
	"github.com/searKing/sole/web/golang/app/configs/values"
)

func Controller() gin.HandlerFunc {
	const IndexTmplName = "web/webapp/WEB-INF/views/index/index.tmpl"
	const IndexHtmlName = "web/webapp/app/modules/index/index.html"

	return func(c *gin.Context) {
		c.Render(http.StatusOK, render.TemplateHTML{
			Name:  "index",
			Files: []string{IndexTmplName},
			Data:  GetIndexTemplateInfo(values.WebApp, IndexHtmlName),
		})
	}
}
