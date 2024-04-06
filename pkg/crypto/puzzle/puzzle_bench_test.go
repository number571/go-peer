package puzzle

import (
	"math/rand"
	"runtime"
	"testing"
	"time"
)

/*
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkPuzzleSequence/20-bit-12            100         463347150 ns/op
--- BENCH: BenchmarkPuzzleSequence/20-bit-12
    puzzle_bench_test.go:68: Timer: 1.284413537s
    puzzle_bench_test.go:68: Timer: 46.334699469s
BenchmarkPuzzleSequence/21-bit-12            100         852326526 ns/op
--- BENCH: BenchmarkPuzzleSequence/21-bit-12
    puzzle_bench_test.go:68: Timer: 2.093841036s
    puzzle_bench_test.go:68: Timer: 1m25.232636438s
BenchmarkPuzzleSequence/22-bit-12            100        1644269096 ns/op
--- BENCH: BenchmarkPuzzleSequence/22-bit-12
    puzzle_bench_test.go:68: Timer: 6.710678809s
    puzzle_bench_test.go:68: Timer: 2m44.426894314s
BenchmarkPuzzleSequence/23-bit-12            100        3301111606 ns/op
--- BENCH: BenchmarkPuzzleSequence/23-bit-12
    puzzle_bench_test.go:68: Timer: 6.587686655s
    puzzle_bench_test.go:68: Timer: 5m30.111137125s
PASS
*/

// go test -bench=BenchmarkPuzzleSequence -benchtime=100x -timeout 99999s
func BenchmarkPuzzleSequence(b *testing.B) {
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
		{
			name:     "23-bit",
			function: NewPoWPuzzle(23).ProofBytes,
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

			b.Logf("Timer (N=%d): %s", b.N, end)
		})
	}
}

/*
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkPuzzleParallel/20-bit-parallel-12                   100         133510360 ns/op
--- BENCH: BenchmarkPuzzleParallel/20-bit-parallel-12
    puzzle_bench_test.go:138: Timer: 454.959149ms
    puzzle_bench_test.go:138: Timer: 13.351011892s
BenchmarkPuzzleParallel/21-bit-parallel-12                   100         259375494 ns/op
--- BENCH: BenchmarkPuzzleParallel/21-bit-parallel-12
    puzzle_bench_test.go:138: Timer: 559.685167ms
    puzzle_bench_test.go:138: Timer: 25.937523604s
BenchmarkPuzzleParallel/22-bit-parallel-12                   100         541665877 ns/op
--- BENCH: BenchmarkPuzzleParallel/22-bit-parallel-12
    puzzle_bench_test.go:138: Timer: 2.213373592s
    puzzle_bench_test.go:138: Timer: 54.166558401s
BenchmarkPuzzleParallel/23-bit-parallel-12                   100        1081891382 ns/op
--- BENCH: BenchmarkPuzzleParallel/23-bit-parallel-12
    puzzle_bench_test.go:138: Timer: 2.219345242s
    puzzle_bench_test.go:138: Timer: 1m48.189016302s
PASS
*/

// go test -bench=BenchmarkPuzzleParallel -benchtime=100x -timeout 99999s
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
		{
			name:     "23-bit-parallel",
			function: NewPoWPuzzle(23).ProofBytes,
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

			b.Logf("Timer (N=%d): %s", b.N, end)
		})
	}
}

func testPseudoRandomBytes(pSeed int) []byte {
	r := rand.New(rand.NewSource(int64(pSeed))) //nolint:gosec
	result := make([]byte, 0, 16)
	for i := 0; i < 16; i++ {
		result = append(result, byte(r.Intn(256)))
	}
	return result
}
