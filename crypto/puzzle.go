package crypto

import (
	"bytes"
	"math"
	"math/big"

	"github.com/number571/go-peer/encoding"
)

var (
	_ Puzzle = &PuzzlePOW{}
)

type PuzzlePOW struct {
	diff uint8
}

func NewPuzzle(diff uint64) Puzzle {
	return &PuzzlePOW{uint8(diff)}
}

// Proof of work by the method of finding the desired hash.
// Hash must start with 'diff' number of zero bits.
func (puzzle *PuzzlePOW) Proof(packHash []byte) uint64 {
	var (
		Target  = big.NewInt(1)
		intHash = big.NewInt(1)
		nonce   = uint64(0)
		hash    []byte
	)
	Target.Lsh(Target, 256-uint(puzzle.diff))
	for nonce < math.MaxUint64 {
		hash = NewSHA256(bytes.Join(
			[][]byte{
				packHash,
				encoding.Uint64ToBytes(nonce),
			},
			[]byte{},
		)).Bytes()
		intHash.SetBytes(hash)
		if intHash.Cmp(Target) == -1 {
			return nonce
		}
		nonce++
	}
	return nonce
}

// Verifies the work of the proof of work function.
func (puzzle *PuzzlePOW) Verify(packHash []byte, nonce uint64) bool {
	intHash := big.NewInt(1)
	Target := big.NewInt(1)
	hash := NewSHA256(bytes.Join(
		[][]byte{
			packHash,
			encoding.Uint64ToBytes(nonce),
		},
		[]byte{},
	)).Bytes()
	intHash.SetBytes(hash)
	Target.Lsh(Target, 256-uint(puzzle.diff))
	return intHash.Cmp(Target) == -1
}
