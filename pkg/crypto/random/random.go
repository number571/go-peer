package random

import (
	"crypto/rand"

	"github.com/number571/go-peer/pkg/encoding"
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
func (p *sStdPRNG) GetBytes(n uint64) []byte {
	slice := make([]byte, n)
	_, err := rand.Read(slice)
	if err != nil {
		// 'return nil' is insecure
		panic(err)
	}
	return slice
}

// Generates a cryptographically strong pseudo-random string.
func (p *sStdPRNG) GetString(n uint64) string {
	return encoding.HexEncode(p.GetBytes(n))[:n]
}

// Generate cryptographically strong pseudo-random uint64 number.
func (p *sStdPRNG) GetUint64() uint64 {
	res := [encoding.CSizeUint64]byte{}
	copy(res[:], p.GetBytes(8))
	return encoding.BytesToUint64(res)
}

// Generate cryptographically strong pseudo-random bool value.
func (p *sStdPRNG) GetBool() bool {
	return p.GetBytes(1)[0]%2 == 0
}
