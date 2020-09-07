// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"net/http"

	http_ "github.com/searKing/golang/go/net/http"
	"github.com/searKing/golang/third_party/github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func (p *Provider) GetRejectInsecureHTTP() negroni.Handler {

	logger := p.Logger().WithField("module", "provider.reject_insecure_http")
	webInfo := p.Proto().GetWeb()

	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		http_.RejectInsecureServerInterceptor(next,
			http_.RejectInsecureWithErrorLog(logrus.New(logger).GetStdLogger()),
			http_.RejectInsecureWithForceHttp(webInfo.GetForceDisableTls() || !webInfo.GetTls().GetEnable()),
			http_.RejectInsecureWithAllowedTlsCidrs(webInfo.GetTls().GetAllowedTlsCidrs()),
			http_.RejectInsecureWithWhitelistedPaths(webInfo.GetTls().GetWhitelistedPaths()))
	})
}
