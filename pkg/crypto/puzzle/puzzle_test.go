package puzzle

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/go-peer/test/_data"
	"golang.org/x/crypto/pbkdf2"
)

const (
	tcN = 2
)

func TestPuzzleDiffSize(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	_ = NewPoWPuzzle(5000, tcN)
}

func TestPuzzle(t *testing.T) {
	t.Parallel()

	var (
		puzzle = NewPoWPuzzle(testutils.TCWorkSize, tcN)
		msg    = []byte("hello, world!")
	)

	hash := hashing.NewSHA256Hasher(msg).ToBytes()
	proof := puzzle.ProofBytes(hash)

	if !puzzle.VerifyBytes(hash, proof) {
		t.Error("proof is invalid")
		return
	}

	if NewPoWPuzzle(25, tcN).VerifyBytes(hash, proof) {
		t.Error("proof 10 with 25 bits is valid?")
		return
	}

	hash[3] = hash[3] ^ 8
	if puzzle.VerifyBytes(hash, proof) {
		t.Error("proof is correct?")
		return
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

/*
SHA256
BenchmarkPuzzle/classic-12                   100         440481055 ns/op
--- BENCH: BenchmarkPuzzle/classic-12
    puzzle_test.go:184: Timer: 1.447560022s
    puzzle_test.go:184: Timer: 44.04798902s
BenchmarkPuzzle/parallel-12                  100         148150797 ns/op
--- BENCH: BenchmarkPuzzle/parallel-12
    puzzle_test.go:184: Timer: 504.5351ms
    puzzle_test.go:184: Timer: 14.81505737s
*/

/*
PBKDF2, tcN = 1
BenchmarkPuzzle/classic-pbkdf2-12            100        1420053053 ns/op
--- BENCH: BenchmarkPuzzle/classic-pbkdf2-12
    puzzle_test.go:184: Timer: 2.680241873s
    puzzle_test.go:184: Timer: 2m22.00527795s
BenchmarkPuzzle/parallel-pbkdf2-12           100         529875241 ns/op
--- BENCH: BenchmarkPuzzle/parallel-pbkdf2-12
    puzzle_test.go:184: Timer: 829.850422ms
    puzzle_test.go:184: Timer: 52.98738265s
*/

/*
PBKDF2, tcN = 2
BenchmarkPuzzle/classic-pbkdf2-12            100        2072725056 ns/op
--- BENCH: BenchmarkPuzzle/classic-pbkdf2-12
    puzzle_test.go:164: Timer: 5.103027543s
    puzzle_test.go:164: Timer: 3m27.272485634s
BenchmarkPuzzle/parallel-pbkdf2-12           100         742923374 ns/op
--- BENCH: BenchmarkPuzzle/parallel-pbkdf2-12
    puzzle_test.go:164: Timer: 1.582913278s
    puzzle_test.go:164: Timer: 1m14.292203667s
*/

/*
PBKDF2, tcN = 4
BenchmarkPuzzle/classic-pbkdf2-12            100        3496014792 ns/op
--- BENCH: BenchmarkPuzzle/classic-pbkdf2-12
    puzzle_test.go:142: Timer: 737.164055ms
    puzzle_test.go:142: Timer: 5m49.601461294s
BenchmarkPuzzle/parallel-pbkdf2-12           100         886905385 ns/op
--- BENCH: BenchmarkPuzzle/parallel-pbkdf2-12
    puzzle_test.go:142: Timer: 173.95526ms
    puzzle_test.go:142: Timer: 1m28.690512668s
*/

// go test -bench=BenchmarkPuzzleParallel -benchtime=100x
func BenchmarkPuzzleParallel(b *testing.B) {
	puzzle := NewPoWPuzzle(20, tcN).(*sPoWPuzzle)

	benchTable := []struct {
		name     string
		function func([]byte) uint64
	}{
		{
			name:     "classic",
			function: puzzle.testClassicProofBytes,
		},
		{
			name:     "parallel",
			function: puzzle.testParallelProofBytes,
		},
		{
			name:     "classic-pbkdf2",
			function: puzzle.testClassicPBKDF2ProofBytes,
		},
		{
			name:     "parallel-pbkdf2",
			function: puzzle.testParallelPBKDF2ProofBytes,
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

func (p *sPoWPuzzle) testClassicProofBytes(packHash []byte) uint64 {
	var (
		target  = big.NewInt(1)
		intHash = big.NewInt(1)
		nonce   = uint64(0)
		hash    []byte
	)
	target.Lsh(target, cHashSizeInBits-uint(p.fDiff))
	for nonce < math.MaxUint64 {
		bNonce := encoding.Uint64ToBytes(nonce)
		hash = hashing.NewSHA256Hasher(bytes.Join(
			[][]byte{
				packHash,
				bNonce[:],
			},
			[]byte{},
		)).ToBytes()
		intHash.SetBytes(hash)
		if intHash.Cmp(target) == -1 {
			return nonce
		}
		nonce++
	}
	return 0
}

func (p *sPoWPuzzle) testParallelProofBytes(packHash []byte) uint64 {
	var (
		parallel = uint64(runtime.GOMAXPROCS(0))
		target   = big.NewInt(1)
	)

	closed := make(chan struct{})
	result := make(chan uint64, parallel)

	target.Lsh(target, cHashSizeInBits-uint(p.fDiff))
	for i := uint64(0); i < parallel; i++ {
		var intHash big.Int
		go func(i uint64) {
			for nonce := i; nonce < math.MaxUint64; nonce += parallel {
				select {
				case <-closed:
					return
				default:
					bNonce := encoding.Uint64ToBytes(nonce)
					hash := hashing.NewSHA256Hasher(bytes.Join(
						[][]byte{
							packHash,
							bNonce[:],
						},
						[]byte{},
					)).ToBytes()
					intHash.SetBytes(hash)
					if intHash.Cmp(target) == -1 {
						result <- nonce
						return
					}
				}
			}
		}(i)
	}

	x := <-result
	close(closed)
	return x
}

func (p *sPoWPuzzle) testClassicPBKDF2ProofBytes(packHash []byte) uint64 {
	var (
		target  = big.NewInt(1)
		intHash = big.NewInt(1)
		nonce   = uint64(0)
		hash    []byte
	)
	target.Lsh(target, cHashSizeInBits-uint(p.fDiff))
	for nonce < math.MaxUint64 {
		bNonce := encoding.Uint64ToBytes(nonce)
		hash = pbkdf2.Key(
			packHash,
			bNonce[:],
			tcN,
			hashing.CSHA256Size,
			sha256.New,
		)
		intHash.SetBytes(hash)
		if intHash.Cmp(target) == -1 {
			return nonce
		}
		nonce++
	}
	return 0
}

func (p *sPoWPuzzle) testParallelPBKDF2ProofBytes(packHash []byte) uint64 {
	var (
		parallel = uint64(runtime.GOMAXPROCS(0))
		target   = big.NewInt(1)
	)

	closed := make(chan struct{})
	result := make(chan uint64, parallel)

	target.Lsh(target, cHashSizeInBits-uint(p.fDiff))
	for i := uint64(0); i < parallel; i++ {
		var intHash big.Int
		go func(i uint64) {
			for nonce := i; nonce < math.MaxUint64; nonce += parallel {
				select {
				case <-closed:
					return
				default:
					bNonce := encoding.Uint64ToBytes(nonce)
					hash := pbkdf2.Key(
						packHash,
						bNonce[:],
						tcN,
						hashing.CSHA256Size,
						sha256.New,
					)
					intHash.SetBytes(hash)
					if intHash.Cmp(target) == -1 {
						result <- nonce
						return
					}
				}
			}
		}(i)
	}

	x := <-result
	close(closed)
	return x
}
