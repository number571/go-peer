package puzzle

import (
	"math/rand"
	"runtime"
	"testing"
	"time"
)

/*
goos: linux
goarch: amd64
pkg: github.com/number571/go-peer/pkg/crypto/puzzle
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkPuzzleSequence/worksize=20-bit-12                   100         506621461 ns/op
--- BENCH: BenchmarkPuzzleSequence/worksize=20-bit-12
    puzzle_bench_test.go:76: Timer (N=1): 1.763921125s
    puzzle_bench_test.go:76: Timer (N=100): 50.662125541s
BenchmarkPuzzleSequence/worksize=21-bit-12                   100         976490983 ns/op
--- BENCH: BenchmarkPuzzleSequence/worksize=21-bit-12
    puzzle_bench_test.go:76: Timer (N=1): 2.12864995s
    puzzle_bench_test.go:76: Timer (N=100): 1m37.649082179s
BenchmarkPuzzleSequence/worksize=22-bit-12                   100        2024441365 ns/op
--- BENCH: BenchmarkPuzzleSequence/worksize=22-bit-12
    puzzle_bench_test.go:76: Timer (N=1): 8.13175731s
    puzzle_bench_test.go:76: Timer (N=100): 3m22.444115905s
BenchmarkPuzzleSequence/worksize=23-bit-12                   100        4016088869 ns/op
--- BENCH: BenchmarkPuzzleSequence/worksize=23-bit-12
    puzzle_bench_test.go:76: Timer (N=1): 7.829686992s
    puzzle_bench_test.go:76: Timer (N=100): 6m41.608867925s
PASS
*/

// go test -bench=BenchmarkPuzzleSequence -benchtime=100x -timeout 99999s
func BenchmarkPuzzleSequence(b *testing.B) {
	benchTable := []struct {
		name     string
		function func([]byte, uint64) uint64
	}{
		{
			name:     "worksize=20-bit",
			function: NewPoWPuzzle(20).ProofBytes,
		},
		{
			name:     "worksize=21-bit",
			function: NewPoWPuzzle(21).ProofBytes,
		},
		{
			name:     "worksize=22-bit",
			function: NewPoWPuzzle(22).ProofBytes,
		},
		{
			name:     "worksize=23-bit",
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
goos: linux
goarch: amd64
pkg: github.com/number571/go-peer/pkg/crypto/puzzle
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkPuzzleParallel/worksize=20-bit-12                   100         204043739 ns/op
--- BENCH: BenchmarkPuzzleParallel/worksize=20-bit-12
    puzzle_bench_test.go:146: Timer (N=1): 412.979992ms
    puzzle_bench_test.go:146: Timer (N=100): 20.404154197s
BenchmarkPuzzleParallel/worksize=21-bit-12                   100         324186197 ns/op
--- BENCH: BenchmarkPuzzleParallel/worksize=21-bit-12
    puzzle_bench_test.go:146: Timer (N=1): 846.078102ms
    puzzle_bench_test.go:146: Timer (N=100): 32.418596955s
BenchmarkPuzzleParallel/worksize=22-bit-12                   100         605896906 ns/op
--- BENCH: BenchmarkPuzzleParallel/worksize=22-bit-12
    puzzle_bench_test.go:146: Timer (N=1): 2.838539071s
    puzzle_bench_test.go:146: Timer (N=100): 1m0.58966026s
BenchmarkPuzzleParallel/worksize=23-bit-12                   100        1287749115 ns/op
--- BENCH: BenchmarkPuzzleParallel/worksize=23-bit-12
    puzzle_bench_test.go:146: Timer (N=1): 2.509988782s
    puzzle_bench_test.go:146: Timer (N=100): 2m8.774882439s
PASS
*/

// go test -bench=BenchmarkPuzzleParallel -benchtime=100x -timeout 99999s
func BenchmarkPuzzleParallel(b *testing.B) {
	benchTable := []struct {
		name     string
		function func([]byte, uint64) uint64
	}{
		{
			name:     "worksize=20-bit",
			function: NewPoWPuzzle(20).ProofBytes,
		},
		{
			name:     "worksize=21-bit",
			function: NewPoWPuzzle(21).ProofBytes,
		},
		{
			name:     "worksize=22-bit",
			function: NewPoWPuzzle(22).ProofBytes,
		},
		{
			name:     "worksize=23-bit",
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
