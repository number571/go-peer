package entropy

import (
	"bytes"

	"github.com/number571/go-peer/pkg/crypto/hashing"
)

var (
	_ IEntropyBooster = &sEntropyBooster{}
)

type sEntropyBooster struct {
	fSalt []byte
	fBits uint64
}

func NewEntropyBooster(bits uint64, salt []byte) IEntropyBooster {
	return &sEntropyBooster{
		fBits: bits,
		fSalt: salt,
	}
}

// Increase entropy by multiple hashing.
func (e *sEntropyBooster) BoostEntropy(data []byte) []byte {
	lim := uint64(1 << e.fBits)

	for i := uint64(0); i < lim; i++ {
		data = hashing.NewSHA256Hasher(bytes.Join(
			[][]byte{
				data,
				e.fSalt,
			},
			[]byte{},
		)).ToBytes()
	}

	return data
}
