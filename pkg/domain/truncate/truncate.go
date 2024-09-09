// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package truncate

import (
	"fmt"

	"github.com/searKing/golang/go/encoding/prettyjson"
	"github.com/searKing/golang/go/strings"
)

const MAX = 1024

func DefaultTruncate(v any) string {
	return Truncate(v, 1024, 10, 2)
}
func Truncate(v any, maxString int, maxBytes int, maxElems int) string {
	data, err := prettyjson.Marshal(v,
		prettyjson.WithEncOptsTruncateString(maxString),
		prettyjson.WithEncOptsTruncateBytes(maxBytes),
		prettyjson.WithEncOptsTruncateSliceOrArray(maxElems),
		prettyjson.WithEncOptsTruncateMap(maxElems),
		prettyjson.WithEncOptsTruncateUrl(true),
		prettyjson.WithEncOptsEscapeHTML(false),
		prettyjson.WithEncOptsOmitEmpty(true),
	)
	if err != nil {
		return strings.Truncate(fmt.Sprintf("%v", v), MAX)
	}
	return string(data)
}
