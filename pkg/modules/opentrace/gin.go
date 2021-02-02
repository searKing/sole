// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opentrace

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	strings_ "github.com/searKing/golang/go/strings"
	"github.com/sirupsen/logrus"
)

func GinHttpTrace(notTrace ...string) gin.HandlerFunc {
	logrus.Info("gin http open trace registered")
	return func(c *gin.Context) {
		if strings_.SliceContains(notTrace, c.Request.URL.Path) {
			c.Next()
			return
		}
		var span opentracing.Span
		opName := c.Request.URL.Path

		// It's very possible that server is fronted by a proxy which could have initiated a trace.
		// If so, we should attempt to join it.
		remoteContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header),
		)

		if err != nil {
			span = opentracing.StartSpan(opName)
		} else {
			span = opentracing.StartSpan(opName, opentracing.ChildOf(remoteContext))
		}
		ext.Component.Set(span, "gin")

		defer span.Finish()

		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))

		c.Next()
		ext.HTTPMethod.Set(span, c.Request.Method)
		statusCode := c.Writer.Status()
		if statusCode >= http.StatusBadRequest {
			ext.Error.Set(span, true)
		}
		ext.HTTPStatusCode.Set(span, uint16(statusCode))
	}
}
