// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package key

import (
	"context"
	"fmt"
	"sync"

	"github.com/searKing/sole/pkg/database/orm"

	"gopkg.in/square/go-jose.v2"
)

type MemoryManager struct {
	keys map[string]*jose.JSONWebKeySet
	sync.RWMutex
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		keys: map[string]*jose.JSONWebKeySet{},
	}
}

func (m *MemoryManager) AddKey(ctx context.Context, setId string, key *jose.JSONWebKey) error {
	m.Lock()
	defer m.Unlock()

	if m.keys[setId] == nil {
		m.keys[setId] = &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
	}
	for _, k := range m.keys[setId].Keys {
		if k.KeyID == key.KeyID {
			return fmt.Errorf(`unable to create key with key_id \"%s\" in set_id \"%s\" because that key_id already exists in the set_id: %w`, key.KeyID, setId, orm.ErrConflict)
		}
	}

	m.keys[setId].Keys = append([]jose.JSONWebKey{*key}, m.keys[setId].Keys...)
	return nil
}

func (m *MemoryManager) AddKeySet(ctx context.Context, setId string, keys *jose.JSONWebKeySet) error {
	for _, key := range keys.Keys {
		if err := m.AddKey(ctx, setId, &key); err != nil {
			return err
		}
	}
	return nil
}

func (m *MemoryManager) GetKey(ctx context.Context, setId, keyId string) (*jose.JSONWebKeySet, error) {
	m.RLock()
	defer m.RUnlock()

	keys, found := m.keys[setId]
	if !found {
		return nil, orm.ErrNotFound
	}

	result := keys.Key(keyId)
	if len(result) == 0 {
		return nil, orm.ErrNotFound
	}

	return &jose.JSONWebKeySet{
		Keys: result,
	}, nil
}

func (m *MemoryManager) GetKeySet(ctx context.Context, setId string) (*jose.JSONWebKeySet, error) {
	m.RLock()
	defer m.RUnlock()

	keys, found := m.keys[setId]
	if !found {
		return nil, orm.ErrNotFound
	}

	if len(keys.Keys) == 0 {
		return nil, orm.ErrNotFound
	}

	return keys, nil
}

func (m *MemoryManager) DeleteKey(ctx context.Context, setId, keyId string) error {
	keys, err := m.GetKeySet(ctx, setId)
	if err != nil {
		return err
	}

	m.Lock()
	defer m.Unlock()

	var results []jose.JSONWebKey
	for _, key := range keys.Keys {
		if key.KeyID != keyId {
			results = append(results, key)
		}
	}
	m.keys[setId].Keys = results

	return nil
}

func (m *MemoryManager) DeleteKeySet(ctx context.Context, setId string) error {
	m.Lock()
	defer m.Unlock()

	delete(m.keys, setId)
	return nil
}
