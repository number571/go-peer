package crypto

import (
	"crypto/rand"

	"github.com/number571/gopeer/encoding"
)

// Generates a cryptographically strong pseudo-random bytes.
func RandBytes(max uint) []byte {
	slice := make([]byte, max)
	_, err := rand.Read(slice)
	if err != nil {
		return nil
	}
	return slice
}

// Generates a cryptographically strong pseudo-random string.
func RandString(max uint) string {
	return encoding.Base64Encode(RandBytes(max))[:max]
}

// Generate cryptographically strong pseudo-random uint64 number.
func RandUint64() uint64 {
	return encoding.BytesToUint64(RandBytes(4))
}
