// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pasta

import (
	"crypto/sha256"

	"github.com/searKing/golang/go/crypto/rand"
)

// GenerateSecret generates random bytes with CharsetAlphaNum
func GenerateSecret(length int) []byte {
	return []byte(rand.StringMath(length))
}

// HashStringSecret hashes the secret for consumption by the AEAD encryption algorithm which expects exactly 32 bytes.
//
// The system secret is being hashed to always match exactly the 32 bytes required by AEAD, even if the secret is long or
// shorter.
func HashStringSecret(secret string) []byte {
	return HashByteSecret([]byte(secret))
}

// HashByteSecret hashes the secret for consumption by the AEAD encryption algorithm which expects exactly 32 bytes.
//
// The system secret is being hashed to always match exactly the 32 bytes required by AEAD, even if the secret is long or
// shorter.
func HashByteSecret(secret []byte) []byte {
	var r [32]byte
	r = sha256.Sum256(secret)
	return r[:]
}
