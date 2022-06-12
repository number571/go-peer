package crypto

import (
	"bytes"
	"math"
	"math/big"

	"github.com/number571/go-peer/encoding"
)

var (
	_ IPuzzle = &sPowPuzzle{}
)

type sPowPuzzle struct {
	fDiff uint8
}

func NewPuzzle(diff uint64) IPuzzle {
	return &sPowPuzzle{uint8(diff)}
}

// Proof of work by the method of finding the desired hash.
// Hash must start with 'diff' number of zero bits.
func (puzzle *sPowPuzzle) Proof(packHash []byte) uint64 {
	var (
		target  = big.NewInt(1)
		intHash = big.NewInt(1)
		nonce   = uint64(0)
		hash    []byte
	)
	target.Lsh(target, sizeInBits(HashSize)-uint(puzzle.fDiff))
	for nonce < math.MaxUint64 {
		hash = NewHasher(bytes.Join(
			[][]byte{
				packHash,
				encoding.Uint64ToBytes(nonce),
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
func (puzzle *sPowPuzzle) Verify(packHash []byte, nonce uint64) bool {
	intHash := big.NewInt(1)
	target := big.NewInt(1)
	hash := NewHasher(bytes.Join(
		[][]byte{
			packHash,
			encoding.Uint64ToBytes(nonce),
		},
		[]byte{},
	)).Bytes()
	intHash.SetBytes(hash)
	target.Lsh(target, sizeInBits(HashSize)-uint(puzzle.fDiff))
	return intHash.Cmp(target) == -1
}

func sizeInBits(n uint) uint {
	return n * 8
}
