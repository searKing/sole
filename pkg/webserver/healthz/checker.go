// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package healthz

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/searKing/golang/go/runtime"
	time_ "github.com/searKing/golang/go/time"
)

// HealthCheck is a named healthz checker.
type HealthCheck interface {
	Name() string
	Check(req *http.Request) error
}

// PingHealthCheck returns true automatically when checked
var PingHealthCheck HealthCheck = ping{}

// ping implements the simplest possible healthz checker.
type ping struct{}

func (ping) Name() string {
	return "ping"
}

// PingHealthCheck is a health check that returns true.
func (ping) Check(_ *http.Request) error {
	return nil
}

// LogHealthCheck returns true if logging is not blocked
var LogHealthCheck HealthCheck = &log{}

type log struct {
	startOnce    sync.Once
	lastVerified atomic.Value
}

func (l *log) Name() string {
	return "log"
}

func (l *log) Check(_ *http.Request) error {
	l.startOnce.Do(func() {
		l.lastVerified.Store(time.Now())
		go time_.Forever(func() {
			defer runtime.NeverPanicButLog.Recover()
			l.lastVerified.Store(time.Now())
		}, time.Minute)
	})

	lastVerified := l.lastVerified.Load().(time.Time)
	if time.Since(lastVerified) < (2 * time.Minute) {
		return nil
	}
	return fmt.Errorf("logging blocked")
}

// healthChecker implements HealthCheck on an arbitrary name and check function.
type healthChecker struct {
	name  string
	check func(r *http.Request) error
}

var _ HealthCheck = &healthChecker{}

func (c *healthChecker) Name() string {
	return c.name
}

func (c *healthChecker) Check(r *http.Request) error {
	return c.check(r)
}
