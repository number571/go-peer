package crypto

import (
	"crypto/rand"

	"github.com/number571/go-peer/encoding"
)

var (
	_ PRNG = &prngT{}
)

type prngT struct {
}

func NewPRNG() PRNG {
	return &prngT{}
}

// Generates a cryptographically strong pseudo-random bytes.
func (r *prngT) Bytes(n uint64) []byte {
	slice := make([]byte, n)
	_, err := rand.Read(slice)
	if err != nil {
		// 'return nil' is insecure
		panic(err)
	}
	return slice
}

// Generates a cryptographically strong pseudo-random string.
func (r *prngT) String(n uint64) string {
	return encoding.Base64Encode(r.Bytes(n))[:n]
}

// Generate cryptographically strong pseudo-random uint64 number.
func (r *prngT) Uint64() uint64 {
	return encoding.BytesToUint64(r.Bytes(8))
}
