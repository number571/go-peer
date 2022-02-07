package crypto

import (
	"crypto/rand"

	"github.com/number571/go-peer/encoding"
)

// Generates a cryptographically strong pseudo-random bytes.
func RandBytes(max uint64) []byte {
	slice := make([]byte, max)
	_, err := rand.Read(slice)
	if err != nil {
		return nil
	}
	return slice
}

// Generates a cryptographically strong pseudo-random string.
func RandString(max uint64) string {
	return encoding.Base64Encode(RandBytes(max))[:max]
}

// Generate cryptographically strong pseudo-random uint64 number.
func RandUint64() uint64 {
	return encoding.BytesToUint64(RandBytes(8))
}
