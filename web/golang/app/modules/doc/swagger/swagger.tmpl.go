// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package swagger

import (
	"html/template"
	"sync"
)

var swaggerIndexTmplCache *template.Template
var swaggerIndexTmplOnce sync.Once

func swaggerIndexTmpl(filenames ...string) *template.Template {
	swaggerIndexTmplOnce.Do(func() {
		// Input: data
		swaggerIndexTmplCache = template.Must(template.New("swaggerTmpl").ParseFiles(filenames...))
	})
	return swaggerIndexTmplCache
}
