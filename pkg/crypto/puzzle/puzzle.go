package puzzle

import (
	"bytes"
	"math"
	"math/big"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
)

const (
	cHashSizeInBits = hashing.CSHA256Size * 8
)

var (
	_ IPuzzle = &sPoWPuzzle{}
)

type sPoWPuzzle struct {
	fDiff uint8
}

func NewPoWPuzzle(pDiff uint64) IPuzzle {
	if pDiff >= math.MaxUint8 {
		panic("diff >= 256")
	}
	return &sPoWPuzzle{
		fDiff: uint8(pDiff),
	}
}

// Proof of work by the method of finding the desired hash.
// Hash must start with 'diff' number of zero bits.
func (p *sPoWPuzzle) ProofBytes(packHash []byte) uint64 {
	var (
		target  = big.NewInt(1)
		intHash = big.NewInt(1)
	)
	target.Lsh(target, cHashSizeInBits-uint(p.fDiff))
	for nonce := uint64(0); nonce < math.MaxUint64; nonce++ {
		bNonce := encoding.Uint64ToBytes(nonce)
		hash := hashing.NewSHA256Hasher(bytes.Join(
			[][]byte{
				packHash,
				bNonce[:],
			},
			[]byte{},
		)).ToBytes()
		intHash.SetBytes(hash)
		if intHash.Cmp(target) == -1 {
			return nonce
		}
	}
	return 0
}

// Verifies the work of the proof of work function.
func (p *sPoWPuzzle) VerifyBytes(packHash []byte, nonce uint64) bool {
	var (
		intHash = big.NewInt(1)
		target  = big.NewInt(1)
	)
	bNonce := encoding.Uint64ToBytes(nonce)
	hash := hashing.NewSHA256Hasher(bytes.Join(
		[][]byte{
			packHash,
			bNonce[:],
		},
		[]byte{},
	)).ToBytes()
	intHash.SetBytes(hash)
	target.Lsh(target, cHashSizeInBits-uint(p.fDiff))
	return intHash.Cmp(target) == -1
}
