package puzzle

import (
	"testing"

	"github.com/number571/go-peer/crypto/hashing"
)

func TestPuzzle(t *testing.T) {
	var (
		puzzle = NewPoWPuzzle(10)
		msg    = []byte("hello, world!")
	)

	hash := hashing.NewSHA256Hasher(msg).Bytes()
	proof := puzzle.Proof(hash)

	if !puzzle.Verify(hash, proof) {
		t.Errorf("proof is invalid")
	}

	if NewPoWPuzzle(25).Verify(hash, proof) {
		t.Errorf("proof 10 with 25 bits is valid?")
	}

	hash[3] = hash[3] ^ 8
	if puzzle.Verify(hash, proof) {
		t.Errorf("proof is correct?")
	}
}
