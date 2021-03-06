// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/searKing/golang/go/net/addr"

	"github.com/searKing/sole/pkg/protobuf"
)

func (p *Provider) HTTPScheme() string {
	return "http"
}

func (p *Provider) GetBackendBindHostPort() string {
	local := p.Proto().GetWeb().GetBindAddr()
	return getHostPort(local.GetHost(), local.GetPort())
}

func (p *Provider) GetBackendAdvertiseHostPort() string {
	addr := p.Proto().GetWeb().GetAdvertiseAddr()
	if addr.GetHost() == "" {
		return p.GetBackendBindHostPort()
	}
	return getHostPort(addr.GetHost(), addr.GetPort())
}

func (p *Provider) GetBackendServeHostPort() string {
	if p.Proto().GetWeb().GetAdvertiseAddr().GetHost() != "" {
		return getHostPort(p.Proto().GetWeb().GetAdvertiseAddr().GetHost(),
			p.Proto().GetWeb().GetAdvertiseAddr().GetPort())
	}
	if p.Proto().GetWeb().GetBindAddr().GetHost() != "" &&
		p.Proto().GetWeb().GetBindAddr().GetHost() != "0.0.0.0" {
		return getHostPort(p.Proto().GetWeb().GetBindAddr().GetHost(),
			p.Proto().GetWeb().GetBindAddr().GetPort())
	}
	resolvers := p.Proto().GetWeb().GetLocalIpResolver()
	ip, err := addr.ServeIP(resolvers.GetNetworks(), resolvers.GetAddresses(),
		protobuf.DurationOrDefault(resolvers.GetTimeout(), 0, "timeout"))
	if err != nil {
		return getHostPort("localhost",
			p.Proto().GetWeb().GetBindAddr().GetPort())
	}
	return getHostPort(ip.String(), p.Proto().GetWeb().GetBindAddr().GetPort())
}

func (p *Provider) ResolveBackendLocalUrl(relativePaths ...string) string {
	return resolveLocalUrl(
		p.HTTPScheme(),
		p.GetBackendServeHostPort(),
		filepath.Join(relativePaths...)).String()
}

func getHostPort(hostname string, port int32) string {
	if strings.HasPrefix(hostname, "unix:") {
		return hostname
	}
	return fmt.Sprintf("%s:%d", hostname, port)
}

func resolveLocalUrl(scheme, hostport, path string) *url.URL {
	u := &url.URL{
		Scheme: scheme,
		Host:   hostport,
		Path:   path,
	}
	if u.Hostname() == "" {
		// use local host
		localHost := "localhost"

		// use local ip
		localIP, err := addr.ListenIP()
		if err == nil && len(localIP) > 0 {
			localHost = localIP.String()
		}
		u.Host = net.JoinHostPort(localHost, u.Port())
	}
	return u
}
