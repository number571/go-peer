package puzzle

import (
	"runtime"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
)

func TestPuzzleDiffSize(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()
	_ = NewPoWPuzzle(5000)
}

func TestPuzzle(t *testing.T) {
	t.Parallel()

	var (
		puzzle = NewPoWPuzzle(10)
		msg    = []byte("hello, world!")
	)

	hash := hashing.NewHasher(msg).ToBytes()
	proof := puzzle.ProofBytes(hash, 1)

	if !puzzle.VerifyBytes(hash, proof) {
		t.Fatal("proof is invalid")
	}

	if puzzle.ProofBytes(hash, 0) != proof {
		t.Fatal("proof is invalid with parallel=[0,1]")
	}

	if NewPoWPuzzle(25).VerifyBytes(hash, proof) {
		t.Fatal("proof 10 with 25 bits is valid?")
	}

	hash[3] ^= 8
	if puzzle.VerifyBytes(hash, proof) {
		t.Fatal("proof is correct?")
	}
}

func TestMultiPuzzle(t *testing.T) {
	t.Parallel()

	parallel := uint64(runtime.GOMAXPROCS(0) + 1) //nolint:gosec
	puzzle := NewPoWPuzzle(4)
	for i := uint64(0); i < 1_000; i++ {
		arr := encoding.Uint64ToBytes(i)
		_ = puzzle.ProofBytes(arr[:], parallel)
	}
}
