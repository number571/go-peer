package hashing

import (
	"testing"
	"time"

	testutils "github.com/number571/go-peer/test/utils"
)

/*
goos: linux
goarch: amd64
pkg: github.com/number571/go-peer/pkg/crypto/hashing
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkHasher-12       1000000               268.5 ns/op
--- BENCH: BenchmarkHasher-12
    hashing_bench_test.go:37: Timer (N=1): 717ns
    hashing_bench_test.go:37: Timer (N=1000000): 268.460938ms
PASS
*/

// go test -bench=BenchmarkHasher -benchtime=1000000x -timeout 99999s
func BenchmarkHasher(b *testing.B) {
	b.StopTimer()
	randomBytes := testutils.PseudoRandomBytes(int(time.Now().UnixNano()))
	b.StartTimer()

	now := time.Now()
	for i := 0; i < b.N; i++ {
		_ = NewHasher(randomBytes).ToBytes()
	}
	end := time.Since(now)

	b.Logf("Timer (N=%d): %s", b.N, end)
}
