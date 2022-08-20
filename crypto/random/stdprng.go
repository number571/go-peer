package random

import (
	"crypto/rand"

	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/settings"
)

var (
	_ IPRNG = &sStdPRNG{}
)

type sStdPRNG struct {
}

func NewStdPRNG() IPRNG {
	return &sStdPRNG{}
}

// Generates a cryptographically strong pseudo-random bytes.
func (r *sStdPRNG) Bytes(n uint64) []byte {
	slice := make([]byte, n)
	_, err := rand.Read(slice)
	if err != nil {
		// 'return nil' is insecure
		panic(err)
	}
	return slice
}

// Generates a cryptographically strong pseudo-random string.
func (r *sStdPRNG) String(n uint64) string {
	return encoding.Base64Encode(r.Bytes(n))[:n]
}

// Generate cryptographically strong pseudo-random uint64 number.
func (r *sStdPRNG) Uint64() uint64 {
	res := [settings.CSizeUint64]byte{}
	copy(res[:], r.Bytes(8))
	return encoding.BytesToUint64(res)
}

// Generate cryptographically strong pseudo-random bool value.
func (r *sStdPRNG) Bool() bool {
	return r.Bytes(1)[0]%2 == 0
}
