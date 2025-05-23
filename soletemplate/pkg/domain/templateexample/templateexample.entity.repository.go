// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package templateexample

import "context"

type Repository interface {
	Example(ctx context.Context, req *ExampleRequest) (resp *ExampleResponse, err error)
}
