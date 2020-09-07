// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package values

import "github.com/searKing/sole/api/protobuf-spec/v1/index"

var (
	Index       = index.Pattern_IndexService_HomePage_0.String() //"/index"
	IndexAsBase = "/"
	IndexAsHtml = index.Pattern_IndexService_HomePage_1.String() // "/index.html"
)
