// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"strings"

	"github.com/gin-gonic/gin"
	gin2 "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
	"github.com/urfave/negroni"
)

func (p *Provider) Negroni() gin.HandlerFunc {
	logger := p.Logger().WithField("module", "provider.middleware.negroni")

	n := negroni.New()
	// use gin.LoggerWithWriter
	//n.Use(negronilogrus.NewMiddlewareFromLogger(h.c.GetLogger(), h.c.Service.Name))

	n.Use(p.PrometheusMetricsManager())
	logger.Infof(`middleware prometheus is loaded`)

	rejectInsecure := !addressIsUnixSocket(p.GetBackendBindHostPort())
	if rejectInsecure {
		n.Use(p.GetRejectInsecureHTTP())
		logger.Infof(`middleware reject insecure http is loaded`)
	}

	if tracer := p.Tracer(); tracer != nil {
		n.Use(tracer)
		logger.Infof(`middleware tracer is loaded`)
	}

	return gin2.UseNegroni(n)

}

func addressIsUnixSocket(address string) bool {
	return strings.HasPrefix(address, "unix:")
}
