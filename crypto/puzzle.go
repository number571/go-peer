package crypto

import (
	"bytes"
	"math"
	"math/big"

	"github.com/number571/gopeer/encoding"
)

var (
	_ Puzzle = &PuzzlePOW{}
)

type PuzzlePOW struct {
	diff uint
}

func NewPuzzle(diff uint) Puzzle {
	return &PuzzlePOW{diff}
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
	Target.Lsh(Target, 256-puzzle.diff)
	for nonce < math.MaxUint64 {
		hash = HashSum(bytes.Join(
			[][]byte{
				packHash,
				encoding.ToBytes(nonce),
			},
			[]byte{},
		))
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
	hash := HashSum(bytes.Join(
		[][]byte{
			packHash,
			encoding.ToBytes(nonce),
		},
		[]byte{},
	))
	intHash.SetBytes(hash)
	Target.Lsh(Target, 256-puzzle.diff)
	return intHash.Cmp(Target) == -1
}
