package puzzle

import (
	"math/rand"
	"runtime"
	"testing"
	"time"
)

/*
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
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
		function func([]byte, uint64) uint64
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
				_ = t.function(randomBytes[i], 1)
			}
			end := time.Since(now)

			b.Log("Timer:", end)
		})
	}
}

/*
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkPuzzleParallel/20-bit-parallel-12                   100         212497876 ns/op
--- BENCH: BenchmarkPuzzleParallel/20-bit-parallel-12
    puzzle_benchmark_test.go:108: Timer: 468.436922ms
    puzzle_benchmark_test.go:108: Timer: 21.249665372s
BenchmarkPuzzleParallel/21-bit-parallel-12                   100         348733248 ns/op
--- BENCH: BenchmarkPuzzleParallel/21-bit-parallel-12
    puzzle_benchmark_test.go:108: Timer: 622.153059ms
    puzzle_benchmark_test.go:108: Timer: 34.873298802s
BenchmarkPuzzleParallel/22-bit-parallel-12                   100         853011126 ns/op
--- BENCH: BenchmarkPuzzleParallel/22-bit-parallel-12
    puzzle_benchmark_test.go:108: Timer: 3.234952925s
    puzzle_benchmark_test.go:108: Timer: 1m25.300982436s
PASS
*/

// go test -bench=BenchmarkPuzzleParallel -benchtime=100x
func BenchmarkPuzzleParallel(b *testing.B) {
	benchTable := []struct {
		name     string
		function func([]byte, uint64) uint64
	}{
		{
			name:     "20-bit-parallel",
			function: NewPoWPuzzle(20).ProofBytes,
		},
		{
			name:     "21-bit-parallel",
			function: NewPoWPuzzle(21).ProofBytes,
		},
		{
			name:     "22-bit-parallel",
			function: NewPoWPuzzle(22).ProofBytes,
		},
	}

	b.ResetTimer()

	for _, t := range benchTable {
		t := t
		b.Run(t.name, func(b *testing.B) {

			b.StopTimer()
			parallel := uint64(runtime.GOMAXPROCS(0))
			randomBytes := make([][]byte, 0, b.N)
			for i := 0; i < b.N; i++ {
				randomBytes = append(randomBytes, testPseudoRandomBytes(i))
			}
			b.StartTimer()

			now := time.Now()
			for i := 0; i < b.N; i++ {
				_ = t.function(randomBytes[i], parallel)
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
