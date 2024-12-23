// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wire

func JoinCleaners(cleanups ...func()) func() {
	return func() {
		for i := range cleanups {
			cleanup := cleanups[len(cleanups)-1-i]
			if cleanup != nil {
				cleanup()
			}
		}
	}
}
