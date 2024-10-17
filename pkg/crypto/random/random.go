package random

import (
	"crypto/rand"
	"strings"

	"github.com/number571/go-peer/pkg/encoding"
)

const (
	charListSize = (1 << 6) // 2^n is a important condition
	charList     = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-`
)

var (
	_ IRandom = &sRandom{}
)

type sRandom struct {
}

func NewRandom() IRandom {
	return &sRandom{}
}

// Generates a cryptographically strong pseudo-random bytes.
func (p *sRandom) GetBytes(n uint64) []byte {
	slice := make([]byte, n)
	if _, err := rand.Read(slice); err != nil {
		panic(err) // 'return nil' is insecure
	}
	return slice
}

// Generates a cryptographically strong pseudo-random string.
/*
	// 1:2
	raw_byte = 2^8 (1 byte)
	hex_byte = 2^4 + 2^4 (2 bytes)

	// 3:4
	raw_byte = 2^8 + 2^8 + 2^8 (3 bytes)
	cur_byte = 2^6 + 2^6 + 2^6 + 2^6 (4 bytes)

	// result
	security[p=128] = random(16) raw_bytes
	security[p=128] = random(32) hex_byte
	security[p=128] = random(21.(3)) cur_byte
*/
func (p *sRandom) GetString(n uint64) string {
	result := strings.Builder{}
	result.Grow(int(n))

	randBytes := p.GetBytes(n)
	for _, b := range randBytes {
		result.WriteByte(charList[b%charListSize])
	}

	return result.String()
}

// Generate cryptographically strong pseudo-random uint64 number.
func (p *sRandom) GetUint64() uint64 {
	res := [encoding.CSizeUint64]byte{}
	copy(res[:], p.GetBytes(encoding.CSizeUint64))
	return encoding.BytesToUint64(res)
}

// Generate cryptographically strong pseudo-random bool value.
func (p *sRandom) GetBool() bool {
	return p.GetBytes(1)[0]%2 == 0
}
