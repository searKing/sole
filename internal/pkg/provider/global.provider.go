// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

var globalProvider *Provider

func SetGlobalProvider(provider *Provider) {
	globalProvider = provider
}

func GlobalProvider() *Provider {
	return globalProvider
}

func InitGlobalProvider(provider *Provider) {
	SetGlobalProvider(provider)
}
