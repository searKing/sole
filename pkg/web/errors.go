// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

func ErrorChain(handlers ...func(err error, handled bool) (error, bool)) func(err error, handled bool) (error, bool) {
	return func(err error, handled bool) (error, bool) {
		for _, h := range handlers {
			if h != nil {
				err, handled = h(err, handled)
			}
		}
		return err, handled
	}
}
