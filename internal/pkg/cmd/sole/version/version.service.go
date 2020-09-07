// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

import (
	"fmt"

	"github.com/searKing/sole/internal/pkg/provider/viper"
)

func service() {
	fmt.Printf("Version:    %s\n", viper.Version)
	fmt.Printf("Git Hash:   %s\n", viper.GitHash)
	fmt.Printf("Build Time: %s\n", viper.BuildTime)
}
