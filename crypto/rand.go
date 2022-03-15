package crypto

import (
	"crypto/rand"

	"github.com/number571/go-peer/encoding"
)

var (
	_ IPRNG = &tPRNG{}
)

type tPRNG struct {
}

func NewPRNG() IPRNG {
	return &tPRNG{}
}

// Generates a cryptographically strong pseudo-random bytes.
func (r *tPRNG) Bytes(n uint64) []byte {
	slice := make([]byte, n)
	_, err := rand.Read(slice)
	if err != nil {
		// 'return nil' is insecure
		panic(err)
	}
	return slice
}

// Generates a cryptographically strong pseudo-random string.
func (r *tPRNG) String(n uint64) string {
	return encoding.Base64Encode(r.Bytes(n))[:n]
}

// Generate cryptographically strong pseudo-random uint64 number.
func (r *tPRNG) Uint64() uint64 {
	return encoding.BytesToUint64(r.Bytes(8))
}
