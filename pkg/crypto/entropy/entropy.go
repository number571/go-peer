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

func NewEntropyBooster(pBits uint64, pSalt []byte) IEntropyBooster {
	return &sEntropyBooster{
		fBits: pBits,
		fSalt: pSalt,
	}
}

// Increase entropy by multiple hashing.
func (p *sEntropyBooster) BoostEntropy(pData []byte) []byte {
	var (
		lim  = uint64(1 << p.fBits)
		data = pData
	)

	for i := uint64(0); i < lim; i++ {
		data = hashing.NewSHA256Hasher(bytes.Join(
			[][]byte{
				data,
				p.fSalt,
			},
			[]byte{},
		)).ToBytes()
	}

	return data
}
