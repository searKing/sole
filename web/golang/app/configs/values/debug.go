// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package values

import "github.com/searKing/sole/api/protobuf-spec/v1/debug"

var (
	DebugPProf  = debug.Pattern_DebugService_PProf_0_For_Gin.String() //"/debug/pprof/*path"
	DebugExpVar = debug.Pattern_DebugService_ExpVar_0.String()        //"/debug/vars"
)
