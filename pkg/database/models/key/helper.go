// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package key

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
)

type KeyGenerator interface {
	Generate(id, use string) (*jose.JSONWebKeySet, error)
}

func CreateKey(ctx context.Context, r Manager, g KeyGenerator, setId string) (*jose.JSONWebKeySet, error) {
	keys, err := g.Generate(uuid.New().String(), "sig")
	if err != nil {
		return nil, errors.WithMessagef(err, "Could not generate JSON Web Key Set \"%s\".", setId)
	}

	for i, k := range keys.Keys {
		k.Use = "sig"
		keys.Keys[i] = k
	}

	if err = r.AddKeySet(ctx, setId, keys); err != nil {
		return nil, errors.WithMessagef(err, "Could not persist JSON Web Key Set \"%s\".", setId)
	}

	return keys, nil
}

func FirstKey(keys []jose.JSONWebKey) *jose.JSONWebKey {
	if len(keys) == 0 {
		return nil
	}
	return &keys[0]
}

func FindKeyByPrefix(set *jose.JSONWebKeySet, prefix string) (key *jose.JSONWebKey, err error) {
	keys, err := FindKeysByPrefix(set, prefix)
	if err != nil {
		return nil, err
	}

	return FirstKey(keys.Keys), nil
}

func FindKeysByPrefix(set *jose.JSONWebKeySet, prefix string) (*jose.JSONWebKeySet, error) {
	keys := new(jose.JSONWebKeySet)

	for _, k := range set.Keys {
		if len(k.KeyID) >= len(prefix)+1 && k.KeyID[:len(prefix)+1] == prefix+":" {
			keys.Keys = append(keys.Keys, k)
		}
	}

	if len(keys.Keys) == 0 {
		return nil, errors.Errorf("Unable to find key with prefix %s in JSON Web Key Set", prefix)
	}

	return keys, nil
}

func PEMBlockForKey(key interface{}) (*pem.Block, error) {
	switch k := key.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}, nil
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}, nil
	default:
		return nil, errors.New("Invalid key type")
	}
}

func Ider(typ, id string) string {
	if id == "" {
		id = uuid.New().String()
	}
	return fmt.Sprintf("%s:%s", typ, id)
}
