// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsn

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/ory/x/sqlcon"
)

// schema://[user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
func ParseDSN(dsn string) (schema string, cfg *mysql.Config, err error) {
	schema = sqlcon.GetDriverName(dsn)
	prefix := fmt.Sprintf("%s://", schema)
	if strings.HasPrefix(dsn, prefix) {
		dsn = strings.TrimPrefix(dsn, prefix)
	}
	// [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	cfg, err = mysql.ParseDSN(dsn)
	return schema, cfg, err
}

// schema://[user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
func GetDSN(schema string, cfg *mysql.Config) string {
	return fmt.Sprintf("%s://%s", schema, cfg.FormatDSN())
}
