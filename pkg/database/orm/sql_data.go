// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orm

import (
	"database/sql"
	"time"
)

type CommonSqlData struct {
	Id        uint      `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	IsDeleted bool         `db:"is_deleted"`
	DeletedAt sql.NullTime `db:"deleted_at"`

	Version uint `db:"version"`
}
