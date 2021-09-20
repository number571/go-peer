package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
)

const (
	HashSize = sha256.Size
)

// Used SHA256.
func SumHash(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

// Used HMAC(SHA256).
func SumHMAC(key, data []byte) []byte {
	hasher := hmac.New(sha256.New, key)
	hasher.Write(data)
	return hasher.Sum(nil)
}
