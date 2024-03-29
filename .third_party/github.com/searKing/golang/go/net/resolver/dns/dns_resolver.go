// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dns

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	rand_ "github.com/searKing/golang/go/math/rand"
	"github.com/searKing/golang/go/net/resolver"
	time_ "github.com/searKing/golang/go/time"
)

// EnableSRVLookups controls whether the DNS resolver attempts to fetch
// addresses from SRV records.  Must not be changed after init time.
var EnableSRVLookups = false

// Globals to stub out in tests.
var newTimer = time.NewTimer

func init() {
	resolver.Register(NewBuilder())
}

const (
	defaultPort       = "443"
	defaultDNSSvrPort = "53"
)

var (
	errMissingAddr = errors.New("dns resolver: missing address")

	// Addresses ending with a colon that is supposed to be the separator
	// between host and port is not allowed.  E.g. "::" is a valid address as
	// it is an IPv6 address (host only) and "[::]:" is invalid as it ends with
	// a colon as the host and port separator
	errEndsWithColon = errors.New("dns resolver: missing port after port-separator colon")
)

var (
	defaultResolver netResolver = net.DefaultResolver
	// To prevent excessive re-resolution, we enforce a rate limit on DNS
	// resolution requests.
	minDNSResRate = 30 * time.Second
)

var customAuthorityDialler = func(authority string) func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		var dialer net.Dialer
		return dialer.DialContext(ctx, network, authority)
	}
}

var customAuthorityResolver = func(authority string) (netResolver, error) {
	host, port, err := parseTarget(authority, defaultDNSSvrPort)
	if err != nil {
		return nil, err
	}

	authorityWithPort := net.JoinHostPort(host, port)

	return &net.Resolver{
		PreferGo: true,
		Dial:     customAuthorityDialler(authorityWithPort),
	}, nil
}

// NewBuilder creates a dnsBuilder which is used to factory DNS resolvers.
func NewBuilder() resolver.Builder {
	return &dnsBuilder{}
}

type dnsBuilder struct{}

// Build creates and starts a DNS resolver that watches the name resolution of the target.
func (b *dnsBuilder) Build(ctx context.Context, target resolver.Target, opts ...resolver.BuildOption) (resolver.Resolver, error) {
	var opt resolver.Build
	opt.ApplyOptions(opts...)
	host, port, err := parseTarget(target.Endpoint, defaultPort)
	if err != nil {
		return nil, err
	}
	cc := opt.ClientConn

	// IP address.
	if ipAddr, ok := formatIP(host); ok {
		addr := []resolver.Address{{Addr: ipAddr + ":" + port}}
		if cc != nil {
			_ = cc.UpdateState(resolver.State{Addresses: addr})
		}
		return deadResolver{
			addrs: addr,
		}, nil
	}

	// DNS address (non-IP).
	ctx, cancel := context.WithCancel(context.Background())
	d := &dnsResolver{
		host:   host,
		port:   port,
		ctx:    ctx,
		cancel: cancel,
		cc:     cc,
		rn:     make(chan struct{}, 1),
	}

	if target.Authority == "" {
		d.resolver = defaultResolver
	} else {
		d.resolver, err = customAuthorityResolver(target.Authority)
		if err != nil {
			return nil, err
		}
	}

	d.wg.Add(1)
	go d.watcher()
	return d, nil
}

// Scheme returns the naming scheme of this resolver builder, which is "dns".
func (b *dnsBuilder) Scheme() string {
	return "dns"
}

type netResolver interface {
	LookupHost(ctx context.Context, host string) (addrs []string, err error)
	LookupSRV(ctx context.Context, service, proto, name string) (cname string, addrs []*net.SRV, err error)
	LookupTXT(ctx context.Context, name string) (txts []string, err error)
}

// deadResolver is a resolver that does nothing.
type deadResolver struct {
	picker resolver.Picker
	addrs  []resolver.Address
}

func (d deadResolver) ResolveOneAddr(ctx context.Context, opts ...resolver.ResolveOneAddrOption) (resolver.Address, error) {
	if len(d.addrs) == 0 {
		return resolver.Address{}, fmt.Errorf("resolve target, but no addr")
	}
	return d.addrs[rand_.Intn(len(d.addrs))], nil
}
func (d deadResolver) ResolveAddr(ctx context.Context, opts ...resolver.ResolveAddrOption) ([]resolver.Address, error) {
	return d.addrs, nil
}
func (deadResolver) ResolveNow(ctx context.Context, opts ...resolver.ResolveNowOption) {}

func (deadResolver) Close() {}

// dnsResolver watches for the name resolution update for a non-IP target.
type dnsResolver struct {
	host     string
	port     string
	resolver netResolver
	ctx      context.Context
	cancel   context.CancelFunc
	cc       resolver.ClientConn
	// rn channel is used by ResolveNow() to force an immediate resolution of the target.
	rn chan struct{}
	// wg is used to enforce Close() to return after the watcher() goroutine has finished.
	// Otherwise, data race will be possible. [Race Example] in dns_resolver_test we
	// replace the real lookup functions with mocked ones to facilitate testing.
	// If Close() doesn't wait for watcher() goroutine finishes, race detector sometimes
	// will warns lookup (READ the lookup function pointers) inside watcher() goroutine
	// has data race with replaceNetFunc (WRITE the lookup function pointers).
	wg sync.WaitGroup
}

func (d *dnsResolver) ResolveOneAddr(ctx context.Context, opts ...resolver.ResolveOneAddrOption) (resolver.Address, error) {
	d.ResolveNow(ctx)
	addrs, err := d.lookupHost()
	if err != nil {
		return resolver.Address{}, err
	}
	if len(addrs) == 0 {
		return resolver.Address{}, fmt.Errorf("resolve target, but no addr")
	}
	return addrs[rand_.Intn(len(addrs))], nil
}

func (d *dnsResolver) ResolveAddr(ctx context.Context, opts ...resolver.ResolveAddrOption) ([]resolver.Address, error) {
	d.ResolveNow(ctx)
	return d.lookupHost()
}

// ResolveNow invoke an immediate resolution of the target that this dnsResolver watches.
func (d *dnsResolver) ResolveNow(ctx context.Context, opts ...resolver.ResolveNowOption) {
	select {
	case d.rn <- struct{}{}:
	default:
	}
}

// Close closes the dnsResolver.
func (d *dnsResolver) Close() {
	d.cancel()
	d.wg.Wait()
}

func (d *dnsResolver) watcher() {
	defer d.wg.Done()

	backoff := time_.NewGrpcExponentialBackOff()
	for {
		addrs, err := d.lookupHost()
		if d.cc != nil {
			if err != nil {
				// Report error to the underlying grpc.ClientConn.
				d.cc.ReportError(err)
			} else {
				err = d.cc.UpdateState(resolver.State{Addresses: addrs})
			}
		}

		var timer *time.Timer
		if err == nil {
			// Success resolving, wait for the next ResolveNow. However, also wait 30 seconds at the very least
			// to prevent constantly re-resolving.
			backoff.Reset()
			timer = newTimer(minDNSResRate)
			select {
			case <-d.ctx.Done():
				timer.Stop()
				return
			case <-d.rn:
			}
		} else {
			// Poll on an error found in DNS Resolver or an error received from ClientConn.
			bc, _ := backoff.NextBackOff()
			timer = newTimer(bc)
		}
		select {
		case <-d.ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

func (d *dnsResolver) lookupSRV(service, proto string) ([]string, error) {
	if !EnableSRVLookups {
		return nil, nil
	}
	var newAddrs []string
	_, srvs, err := d.resolver.LookupSRV(d.ctx, service, proto, d.host)
	if err != nil {
		err = handleDNSError(err, "SRV") // may become nil
		return nil, err
	}
	for _, s := range srvs {
		lbAddrs, err := d.resolver.LookupHost(d.ctx, s.Target)
		if err != nil {
			err = handleDNSError(err, "A") // may become nil
			if err == nil {
				// If there are other SRV records, look them up and ignore this
				// one that does not exist.
				continue
			}
			return nil, err
		}
		for _, a := range lbAddrs {
			ip, ok := formatIP(a)
			if !ok {
				return nil, fmt.Errorf("dns: error parsing A record IP address %v", a)
			}
			addr := ip + ":" + strconv.Itoa(int(s.Port))
			newAddrs = append(newAddrs, addr)
		}
	}
	return newAddrs, nil
}

var filterError = func(err error) error {
	if dnsErr, ok := err.(*net.DNSError); ok && !dnsErr.IsTimeout && !dnsErr.IsTemporary {
		// Timeouts and temporary errors should be communicated to gRPC to
		// attempt another DNS query (with backoff).  Other errors should be
		// suppressed (they may represent the absence of a TXT record).
		return nil
	}
	return err
}

func handleDNSError(err error, lookupType string) error {
	err = filterError(err)
	if err != nil {
		err = fmt.Errorf("dns: %v record lookup error: %w", lookupType, err)
		return err
	}
	return nil
}

func (d *dnsResolver) lookupHost() ([]resolver.Address, error) {
	var newAddrs []resolver.Address
	addrs, err := d.resolver.LookupHost(d.ctx, d.host)
	if err != nil {
		err = handleDNSError(err, "A")
		return nil, err
	}
	for _, a := range addrs {
		ip, ok := formatIP(a)
		if !ok {
			return nil, fmt.Errorf("dns: error parsing A record IP address %v", a)
		}
		addr := ip + ":" + d.port
		newAddrs = append(newAddrs, resolver.Address{Addr: addr})
	}
	return newAddrs, nil
}

// formatIP returns ok = false if addr is not a valid textual representation of an IP address.
// If addr is an IPv4 address, return the addr and ok = true.
// If addr is an IPv6 address, return the addr enclosed in square brackets and ok = true.
func formatIP(addr string) (addrIP string, ok bool) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return "", false
	}
	if ip.To4() != nil {
		return addr, true
	}
	return "[" + addr + "]", true
}

// parseTarget takes the user input target string and default port, returns formatted host and port info.
// If target doesn't specify a port, set the port to be the defaultPort.
// If target is in IPv6 format and host-name is enclosed in square brackets, brackets
// are stripped when setting the host.
// examples:
// target: "www.google.com" defaultPort: "443" returns host: "www.google.com", port: "443"
// target: "ipv4-host:80" defaultPort: "443" returns host: "ipv4-host", port: "80"
// target: "[ipv6-host]" defaultPort: "443" returns host: "ipv6-host", port: "443"
// target: ":80" defaultPort: "443" returns host: "localhost", port: "80"
func parseTarget(target, defaultPort string) (host, port string, err error) {
	if target == "" {
		return "", "", errMissingAddr
	}
	if ip := net.ParseIP(target); ip != nil {
		// target is an IPv4 or IPv6(without brackets) address
		return target, defaultPort, nil
	}
	if host, port, err = net.SplitHostPort(target); err == nil {
		if port == "" {
			// If the port field is empty (target ends with colon), e.g. "[::1]:", this is an error.
			return "", "", errEndsWithColon
		}
		// target has port, i.e ipv4-host:port, [ipv6-host]:port, host-name:port
		if host == "" {
			// Keep consistent with net.Dial(): If the host is empty, as in ":80", the local system is assumed.
			host = "localhost"
		}
		return host, port, nil
	}
	if host, port, err = net.SplitHostPort(target + ":" + defaultPort); err == nil {
		// target doesn't have port
		return host, port, nil
	}
	return "", "", fmt.Errorf("invalid target address %v, error info: %v", target, err)
}
