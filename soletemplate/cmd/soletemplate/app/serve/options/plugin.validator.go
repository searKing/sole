// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import "github.com/go-playground/validator/v10"

func NewValidator() *validator.Validate { return validator.New() }
