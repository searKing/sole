// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opentrace

//go:generate go-enum -type Type -trimprefix=Type
type Type int

const (
	TypeJeager Type = iota
	TypeZipkin
	TypeButt
)
