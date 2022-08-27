package puzzle

import (
	"bytes"
	"math"
	"math/big"

	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/encoding"
)

var (
	_ IPuzzle = &sPoWPuzzle{}
)

type sPoWPuzzle struct {
	fDiff uint8
}

func NewPoWPuzzle(diff uint64) IPuzzle {
	return &sPoWPuzzle{uint8(diff)}
}

// Proof of work by the method of finding the desired hash.
// Hash must start with 'diff' number of zero bits.
func (puzzle *sPoWPuzzle) Proof(packHash []byte) uint64 {
	var (
		target  = big.NewInt(1)
		intHash = big.NewInt(1)
		nonce   = uint64(0)
		hash    []byte
	)
	target.Lsh(target, hashSizeInBits()-uint(puzzle.fDiff))
	for nonce < math.MaxUint64 {
		bNonce := encoding.Uint64ToBytes(nonce)
		hash = hashing.NewSHA256Hasher(bytes.Join(
			[][]byte{
				packHash,
				bNonce[:],
			},
			[]byte{},
		)).Bytes()
		intHash.SetBytes(hash)
		if intHash.Cmp(target) == -1 {
			return nonce
		}
		nonce++
	}
	return nonce
}

// Verifies the work of the proof of work function.
func (puzzle *sPoWPuzzle) Verify(packHash []byte, nonce uint64) bool {
	intHash := big.NewInt(1)
	target := big.NewInt(1)
	bNonce := encoding.Uint64ToBytes(nonce)
	hash := hashing.NewSHA256Hasher(bytes.Join(
		[][]byte{
			packHash,
			bNonce[:],
		},
		[]byte{},
	)).Bytes()
	intHash.SetBytes(hash)
	target.Lsh(target, hashSizeInBits()-uint(puzzle.fDiff))
	return intHash.Cmp(target) == -1
}

func hashSizeInBits() uint {
	return uint(hashing.CSHA256Size * 8)
}
