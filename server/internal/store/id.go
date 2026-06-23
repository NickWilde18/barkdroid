package store

import (
	"crypto/rand"
	"encoding/hex"
)

func newID() string {
	return randomHex(16)
}

func newKey() string {
	return randomHex(4) // 8-char key, like Bark
}

func randomHex(bytes int) string {
	b := make([]byte, bytes)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand failed: " + err.Error())
	}
	return hex.EncodeToString(b)
}
