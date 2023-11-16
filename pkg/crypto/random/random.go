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

func init() {
	if len(charList) != charListSize {
		panic("len(charList) != charListSize")
	}
}

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
	result := strings.Builder{}
	result.Grow(int(n))

	randBytes := p.GetBytes(n)
	for _, b := range randBytes {
		result.WriteByte(charList[b%charListSize])
	}

	return result.String()
}

// Generate cryptographically strong pseudo-random uint64 number.
func (p *sStdPRNG) GetUint64() uint64 {
	res := [encoding.CSizeUint64]byte{}
	copy(res[:], p.GetBytes(encoding.CSizeUint64))
	return encoding.BytesToUint64(res)
}

// Generate cryptographically strong pseudo-random bool value.
func (p *sStdPRNG) GetBool() bool {
	return p.GetBytes(1)[0]%2 == 0
}
