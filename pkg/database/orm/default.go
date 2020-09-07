// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orm

const (
	// Pagination
	DefaultQueryLimit  = 10
	DefaultQueryOffset = 0
)

func FormatLimit(limit int) int {
	if limit <= 0 {
		limit = DefaultQueryLimit
	}
	return limit
}

func FormatOffset(offset int) int {
	if offset <= 0 {
		offset = DefaultQueryOffset
	}
	return offset
}

func FormatPagination(LimitPerPage, PageSeq int) (limit, offset int) {
	limit = LimitPerPage
	offset = PageSeq * LimitPerPage
	if limit <= 0 {
		limit = DefaultQueryLimit
	}
	if offset <= 0 {
		offset = DefaultQueryOffset
	}
	return limit, offset
}
