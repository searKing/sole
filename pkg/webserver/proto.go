// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/searKing/golang/go/net/addr"
	strings_ "github.com/searKing/golang/go/strings"
	"github.com/searKing/sole/pkg/protobuf"
)

func (web *Web) HTTPScheme() string {
	if web.GetForceDisableTls() {
		return "http"
	}
	return "https"
}

func (web *Web) ResolveLocalIp() string {
	resolvers := web.GetLocalIpResolver()
	ip, err := addr.ServeIP(resolvers.GetNetworks(), resolvers.GetAddresses(),
		protobuf.DurationOrDefault(resolvers.GetTimeout(), 0, "timeout"))
	if err != nil {
		return "localhost"
	}
	return ip.String()
}

// GetBackendBindHostPort returns a address to listen.
func (web *Web) GetBackendBindHostPort() string {
	local := web.GetBindAddr()
	return getHostPort(local.GetHost(), local.GetPort())
}

// GetBackendAdvertiseHostPort returns a address to expose with domain, if not set, use host instead.
func (web *Web) GetBackendAdvertiseHostPort() string {
	extern := web.GetAdvertiseAddr()
	host := strings_.ValueOrDefault(extern.GetDomains()...)
	if host == "" {
		host = web.GetAdvertiseAddr().GetHost()
	}
	if host == "" {
		return web.GetBackendBindHostPort()
	}
	return getHostPort(host, extern.GetPort())
}

// GetBackendServeHostPort returns a address to expose without domain, if not set, use resolver to resolve a ip
func (web *Web) GetBackendServeHostPort() string {
	host := web.GetAdvertiseAddr().GetHost()
	if host != "" {
		return getHostPort(host, web.GetAdvertiseAddr().GetPort())
	}

	host = web.GetBindAddr().GetHost()
	if host != "" && host != "0.0.0.0" {
		return getHostPort(host, web.GetBindAddr().GetPort())
	}
	return getHostPort(web.ResolveLocalIp(), web.GetBindAddr().GetPort())
}

func (web *Web) ResolveBackendLocalUrl(relativePaths ...string) string {
	return resolveLocalUrl(
		web.HTTPScheme(),
		web.GetBackendServeHostPort(),
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
