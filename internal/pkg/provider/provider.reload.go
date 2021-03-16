// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"context"
	"time"

	time_ "github.com/searKing/golang/go/time"
)

const (
	DefaultTimeout = time.Minute
)

func (p *Provider) ReloadForever() {
	p.reloadOnce.Do(func() {
		func() {
			time_.NonSlidingUntil(p.Context(), func(ctx context.Context) {
				providerReloads.WithLabelValues(p.proto.String()).Inc()
			}, DefaultTimeout)
		}()
	})
}
