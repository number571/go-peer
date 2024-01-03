package puzzle

import (
	"math/rand"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestPuzzleDiffSize(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	_ = NewPoWPuzzle(5000)
}

func TestPuzzle(t *testing.T) {
	t.Parallel()

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

/*
BenchmarkPuzzle/20-bit-12                    100         451693442 ns/op
--- BENCH: BenchmarkPuzzle/20-bit-12
    puzzle_test.go:90: Timer: 1.267312821s
    puzzle_test.go:90: Timer: 45.169327733s
BenchmarkPuzzle/21-bit-12                    100         854375738 ns/op
--- BENCH: BenchmarkPuzzle/21-bit-12
    puzzle_test.go:90: Timer: 1.881050349s
    puzzle_test.go:90: Timer: 1m25.437556403s
BenchmarkPuzzle/22-bit-12                    100        1746838647 ns/op
--- BENCH: BenchmarkPuzzle/22-bit-12
    puzzle_test.go:90: Timer: 7.615185783s
    puzzle_test.go:90: Timer: 2m54.6838481s
*/

// go test -bench=BenchmarkPuzzle -benchtime=100x
func BenchmarkPuzzle(b *testing.B) {
	benchTable := []struct {
		name     string
		function func([]byte) uint64
	}{
		{
			name:     "20-bit",
			function: NewPoWPuzzle(20).ProofBytes,
		},
		{
			name:     "21-bit",
			function: NewPoWPuzzle(21).ProofBytes,
		},
		{
			name:     "22-bit",
			function: NewPoWPuzzle(22).ProofBytes,
		},
	}

	b.ResetTimer()

	for _, t := range benchTable {
		t := t
		b.Run(t.name, func(b *testing.B) {
			b.StopTimer()
			randomBytes := make([][]byte, 0, b.N)
			for i := 0; i < b.N; i++ {
				randomBytes = append(randomBytes, testPseudoRandomBytes(i))
			}
			b.StartTimer()

			now := time.Now()
			for i := 0; i < b.N; i++ {
				_ = t.function(randomBytes[i])
			}
			end := time.Since(now)

			b.Log("Timer:", end)
		})
	}
}

func testPseudoRandomBytes(pSeed int) []byte {
	r := rand.New(rand.NewSource(int64(pSeed)))
	result := make([]byte, 0, 16)
	for i := 0; i < 16; i++ {
		result = append(result, byte(r.Intn(256)))
	}
	return result
}
