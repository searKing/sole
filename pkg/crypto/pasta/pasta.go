// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pasta

import (
	"github.com/searKing/golang/third_party/github.com/gtank/cryptopasta"
)

type KeyProvider interface {
	GetSystemSecret() []byte           // primary key
	GetRotatedSystemSecrets() [][]byte // multi keys for multi encrypt|decrypt
}

// Pasta copy & paste-friendly golang crypto
type Pasta struct {
	keyProvider KeyProvider
}

func New(c KeyProvider) *Pasta {
	return &Pasta{keyProvider: c}
}

func NewFromKey(systemSecret []byte, rotatedSystemSecrets [][]byte) *Pasta {
	return &Pasta{keyProvider: &innerKeyCipherProvider{
		systemSecret:         systemSecret,
		rotatedSystemSecrets: rotatedSystemSecrets,
	}}
}

func (c *Pasta) keys() [][]byte {
	return append([][]byte{c.keyProvider.GetSystemSecret()}, c.keyProvider.GetRotatedSystemSecrets()...)
}

func (c *Pasta) Encrypt(plaintext []byte) (string, error) {
	return cryptopasta.Encrypt(plaintext, c.keys()...)
}

func (c *Pasta) Decrypt(ciphertext string) (p []byte, err error) {
	return cryptopasta.Decrypt(ciphertext, c.keys()...)
}

type innerKeyCipherProvider struct {
	systemSecret         []byte
	rotatedSystemSecrets [][]byte
}

func (c *innerKeyCipherProvider) GetRotatedSystemSecrets() [][]byte {
	return c.rotatedSystemSecrets
}

func (c *innerKeyCipherProvider) GetSystemSecret() []byte {
	return c.systemSecret
}
