// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	gincors "github.com/rs/cors/wrapper/gin"
	gin2 "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
)

func (p *Provider) GetCORS() gin.HandlerFunc {
	logger := p.Logger().WithField("module", "provider.web.cors")
	corsInfo := p.Proto().GetWeb().GetCors()
	if corsInfo == nil || !corsInfo.Enable {
		return gin2.NopHandlerFunc
	}

	if !corsInfo.GetUseConditional() {
		logger.Info("Enabled Unconditional CORS")
		return func(ctx *gin.Context) {
			ctx.Header("Access-Control-Allow-Origin", "*")
		}
	}

	maxWait, err := ptypes.Duration(corsInfo.GetMaxAge())
	if err != nil {
		maxWait = 0 * time.Second
		logger.WithField("max_age", corsInfo.GetMaxAge()).
			WithError(errors.WithStack(err)).
			Warnf("malformed max_age, use %s instead", maxWait)
	}
	logger.Info("Enabled Conditional CORS")

	return gincors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		Debug:            corsInfo.GetDebug(),
		MaxAge:           int(maxWait.Seconds()),
	})
	//return gincors.AllowAll()
}
