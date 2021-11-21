package crypto

import (
	"crypto/rand"

	"github.com/number571/gopeer/encoding"
)

// Generates a cryptographically strong pseudo-random sequence.
func RandBytes(max uint) []byte {
	slice := make([]byte, max)
	_, err := rand.Read(slice)
	if err != nil {
		return nil
	}
	return slice
}

// Generate cryptographically strong pseudo-random uint64 number.
func RandUint64() uint64 {
	return encoding.BytesToUint64(RandBytes(4))
}
