// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orm

import (
	"errors"
)

var ErrNotFound = errors.New("unable to locate the requested resource")
var ErrConflict = errors.New("unable to process the requested resource because of conflict in the current state")
var ErrUnImplemented = errors.New("unimplemented method")
