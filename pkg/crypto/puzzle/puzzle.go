package puzzle

import (
	"crypto/sha256"
	"math"
	"math/big"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	"golang.org/x/crypto/pbkdf2"
)

const (
	cHashSizeInBits = hashing.CSHA256Size * 8
)

var (
	_ IPuzzle = &sPoWPuzzle{}
)

type sPoWPuzzle struct {
	fDiff uint8
	fIter uint64
}

func NewPoWPuzzle(pDiff, pIter uint64) IPuzzle {
	if pDiff >= math.MaxUint8 {
		panic("diff >= 256")
	}
	return &sPoWPuzzle{
		fDiff: uint8(pDiff),
		fIter: pIter,
	}
}

// Proof of work by the method of finding the desired hash.
// Hash must start with 'diff' number of zero bits.
func (p *sPoWPuzzle) ProofBytes(packHash []byte) uint64 {
	var (
		target  = big.NewInt(1)
		intHash = big.NewInt(1)
		nonce   = uint64(0)
	)
	target.Lsh(target, cHashSizeInBits-uint(p.fDiff))
	for nonce < math.MaxUint64 {
		bNonce := encoding.Uint64ToBytes(nonce)
		hash := pbkdf2.Key(
			packHash,
			bNonce[:],
			int(p.fIter),
			hashing.CSHA256Size,
			sha256.New,
		)
		intHash.SetBytes(hash)
		if intHash.Cmp(target) == -1 {
			return nonce
		}
		nonce++
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
	hash := pbkdf2.Key(
		packHash,
		bNonce[:],
		int(p.fIter),
		hashing.CSHA256Size,
		sha256.New,
	)
	intHash.SetBytes(hash)
	target.Lsh(target, cHashSizeInBits-uint(p.fDiff))
	return intHash.Cmp(target) == -1
}
