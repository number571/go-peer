package puzzle

import (
	"testing"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestPuzzle(t *testing.T) {
	var (
		puzzle = NewPoWPuzzle(testutils.TCWorkSize)
		msg    = []byte("hello, world!")
	)

	hash := hashing.NewSHA256Hasher(msg).ToBytes()
	proof := puzzle.ProofBytes(hash)

	if !puzzle.VerifyBytes(hash, proof) {
		t.Error("proof is invalid")
		return
	}

	if NewPoWPuzzle(25).VerifyBytes(hash, proof) {
		t.Error("proof 10 with 25 bits is valid?")
		return
	}

	hash[3] = hash[3] ^ 8
	if puzzle.VerifyBytes(hash, proof) {
		t.Error("proof is correct?")
		return
	}
}
