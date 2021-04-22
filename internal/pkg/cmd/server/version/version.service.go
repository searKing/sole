// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

import (
	"fmt"

	"github.com/searKing/sole/pkg/appinfo"
)

func Run() {
	fmt.Printf("Version:    %s\n", appinfo.Version)
	fmt.Printf("Git Hash:   %s\n", appinfo.GitHash)
	fmt.Printf("Build Time: %s\n", appinfo.BuildTime)
}
