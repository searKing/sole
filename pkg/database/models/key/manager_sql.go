// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package key

import (
	"context"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/ory/x/sqlcon"
	"github.com/pkg/errors"
	"github.com/searKing/golang/go/util/object"
	"github.com/searKing/sole/pkg/crypto/pasta"
	"github.com/searKing/sole/pkg/database/orm"
	"gopkg.in/square/go-jose.v2"
)

type Provider interface {
	KeyCipher() *pasta.Pasta
}

type SQLManager struct {
	db       *sqlx.DB
	provider Provider
}

const (
	version      = 0
	addKeySet    = `INSERT INTO sole_key (set_id, key_id, key_data, version) VALUES (:set_id, :key_id, :key_data, :version)`
	queryKey     = `SELECT * FROM sole_key WHERE set_id=? AND key_id=? ORDER BY created_at DESC`
	queryKeySet  = `SELECT * FROM sole_key WHERE set_id=?  ORDER BY created_at DESC`
	deleteKey    = `DELETE FROM sole_key WHERE set_id=? AND key_id=?`
	deleteKeySet = `DELETE FROM sole_key WHERE set_id=?`
)

type sqlData struct {
	SetId   string `db:"set_id"`
	KeyId   string `db:"key_id"`
	KeyData string `db:"key_data"`
	orm.CommonSqlData
	//Id        int       `db:"id"`
	//CreatedAt time.Time `db:"created_at"`
	//UpdatedAt time.Time `db:"updated_at"`
	//
	//IsDeleted bool       `db:"is_deleted"`
	//DeletedAt *time.Time `db:"deleted_at"`
	//
	//Version int `db:"version"`
}

// TableName return table name
func (*sqlData) TableName() string {
	return `sole_key`
}

func NewSQLManager(db *sqlx.DB, provider Provider) *SQLManager {
	return &SQLManager{
		db:       db,
		provider: provider,
	}
}

func (*SQLManager) MigrateTableName() string {
	return `sole_key_migration`
}

func (m *SQLManager) AddKey(ctx context.Context, setId string, key *jose.JSONWebKey) error {
	object.RequireNonNull(m.db)
	object.RequireNonNull(m.provider)
	out, err := json.Marshal(key)
	if err != nil {
		return err
	}

	encrypted, err := m.provider.KeyCipher().Encrypt(out)
	if err != nil {
		return err
	}

	if _, err = m.db.NamedExecContext(ctx, addKeySet, &sqlData{
		SetId:   setId,
		KeyId:   key.KeyID,
		KeyData: encrypted,
		CommonSqlData: orm.CommonSqlData{
			Version: version,
		},
	}); err != nil {
		return sqlcon.HandleError(err)
	}
	return nil
}

func (m *SQLManager) AddKeySet(ctx context.Context, setId string, keys *jose.JSONWebKeySet) error {
	object.RequireNonNull(m.db)
	object.RequireNonNull(m.provider)
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err := m.addKeySet(ctx, tx, m.provider.KeyCipher(), setId, keys); err != nil {
		if re := tx.Rollback(); re != nil {
			return errors.Wrap(err, re.Error())
		}
		return sqlcon.HandleError(err)
	}

	if err := tx.Commit(); err != nil {
		if re := tx.Rollback(); re != nil {
			return errors.Wrap(err, re.Error())
		}
		return sqlcon.HandleError(err)
	}
	return nil
}

func (m *SQLManager) addKeySet(ctx context.Context, tx *sqlx.Tx, cipher *pasta.Pasta, setId string, keys *jose.JSONWebKeySet) error {
	for _, key := range keys.Keys {
		out, err := json.Marshal(key)
		if err != nil {
			return err
		}

		encrypted, err := cipher.Encrypt(out)
		if err != nil {
			return err
		}

		if _, err = tx.NamedExecContext(ctx, addKeySet, &sqlData{
			SetId:   setId,
			KeyId:   key.KeyID,
			KeyData: encrypted,
			CommonSqlData: orm.CommonSqlData{
				Version: version,
			},
		}); err != nil {
			return sqlcon.HandleError(err)
		}
	}

	return nil
}

func (m *SQLManager) GetKey(ctx context.Context, setId, keyId string) (*jose.JSONWebKeySet, error) {
	object.RequireNonNull(m.db)
	object.RequireNonNull(m.provider)
	var d sqlData
	if err := m.db.GetContext(ctx, &d, m.db.Rebind(queryKey), setId, keyId); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	key, err := m.provider.KeyCipher().Decrypt(d.KeyData)
	if err != nil {
		return nil, err
	}

	var c jose.JSONWebKey
	if err := json.Unmarshal(key, &c); err != nil {
		return nil, err
	}

	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{c},
	}, nil
}

func (m *SQLManager) GetKeySet(ctx context.Context, setId string) (*jose.JSONWebKeySet, error) {
	object.RequireNonNull(m.db)
	object.RequireNonNull(m.provider)
	var ds []sqlData
	if err := m.db.SelectContext(ctx, &ds, m.db.Rebind(queryKeySet), setId); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	if len(ds) == 0 {
		return nil, orm.ErrNotFound
	}

	keys := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
	for _, d := range ds {
		key, err := m.provider.KeyCipher().Decrypt(d.KeyData)
		if err != nil {
			return nil, err
		}

		var c jose.JSONWebKey
		if err := json.Unmarshal(key, &c); err != nil {
			return nil, err
		}
		keys.Keys = append(keys.Keys, c)
	}

	if len(keys.Keys) == 0 {
		return nil, orm.ErrNotFound
	}

	return keys, nil
}

func (m *SQLManager) DeleteKey(ctx context.Context, setId, keyId string) error {
	object.RequireNonNull(m.db)
	if _, err := m.db.ExecContext(ctx, m.db.Rebind(deleteKey), setId, keyId); err != nil {
		return sqlcon.HandleError(err)
	}
	return nil
}

func (m *SQLManager) DeleteKeySet(ctx context.Context, setId string) error {
	object.RequireNonNull(m.db)
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err := m.deleteKeySet(ctx, tx, setId); err != nil {
		if re := tx.Rollback(); re != nil {
			return errors.Wrap(err, re.Error())
		}
		return sqlcon.HandleError(err)
	}

	if err := tx.Commit(); err != nil {
		if re := tx.Rollback(); re != nil {
			return errors.Wrap(err, re.Error())
		}
		return sqlcon.HandleError(err)
	}
	return nil
}

func (m *SQLManager) deleteKeySet(ctx context.Context, tx *sqlx.Tx, setId string) error {
	if _, err := tx.ExecContext(ctx, m.db.Rebind(deleteKeySet), setId); err != nil {
		return sqlcon.HandleError(err)
	}
	return nil
}
