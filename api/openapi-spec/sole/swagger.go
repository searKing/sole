// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sole

import _ "embed"

//go:embed swagger.json
var SwaggerUIJson string

//go:embed swagger.yaml
var SwaggerUIYaml string
