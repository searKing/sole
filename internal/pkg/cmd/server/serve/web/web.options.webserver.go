// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"

	"github.com/searKing/sole/pkg/net/cors"
	"github.com/searKing/sole/web/golang"
)

func (s *ServerRunOptions) completeWebServer() error {
	s.WebServerOptions.BindAddress = s.Provider.GetBackendBindHostPort()
	s.WebServerOptions.ExternalAddress = s.Provider.GetBackendServeHostPort()
	s.WebServerOptions.AddWebHandler(golang.NewHandler())

	{
		corsConfig := cors.NewConfig()
		corsInfo := s.Provider.Proto().GetWeb().GetCors()
		if corsInfo != nil {
			if corsInfo.Enable {
				maxAge, err := ptypes.Duration(corsInfo.GetMaxAge())
				if err != nil {
					maxAge = 0 * time.Second
					logrus.WithField("max_age", corsInfo.GetMaxAge()).
						WithError(err).
						Warnf("malformed max_age, use %s instead", maxAge)
				}
				corsConfig.UseConditional = corsInfo.GetUseConditional()
				corsConfig.AllowedOrigins = corsInfo.GetAllowedOrigins()
				corsConfig.AllowedMethods = corsInfo.GetAllowedOrigins()
				corsConfig.AllowedHeaders = corsInfo.GetAllowedHeaders()
				corsConfig.ExposedHeaders = corsInfo.GetExposedHeaders()

				corsConfig.MaxAge = maxAge
				corsConfig.AllowCredentials = corsInfo.GetAllowCredentials()
			} else {
				corsConfig = nil
			}
		}
		s.WebServerOptions.CORS = corsConfig
	}
	{
		consulInfo := s.Provider.Proto().GetConsul()
		if consulInfo.GetEnabled() {
			backend, err := s.ServiceRegistry.Complete().New()
			if err != nil {
				logrus.WithError(err).Errorf("build service registry backend")
				return err
			}
			s.WebServerOptions.ServiceRegistryBackend = backend
		}
	}
	return nil
}
