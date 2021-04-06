// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

func (s *ServerRunOptions) completeServiceResgistry() error {
	s.ServiceRegistry.ServiceName = s.Provider.Proto().GetService().GetName()
	s.ServiceRegistry.ServiceAddress = s.Provider.GetBackendServeHostPort()
	s.ServiceRegistry.ConsulAddress = s.Provider.Proto().GetConsul().GetAddress()
	s.ServiceRegistry.HealthCheckUrl = "/healthz"
	return nil
}
