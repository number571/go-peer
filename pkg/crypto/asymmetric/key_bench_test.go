package asymmetric

import (
	"testing"
	"time"
)

/*
goos: linux
goarch: amd64
pkg: github.com/number571/go-peer/pkg/crypto/asymmetric
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkKey/mlkem=768,mldsa=65-12                 10000            171410 ns/op
--- BENCH: BenchmarkKey/mlkem=768,mldsa=65-12
    key_bench_test.go:26: Timer_New(N=1): 211.39µs
    key_bench_test.go:26: Timer_New(N=10000): 1.714080663s
PASS
*/

// go test -bench=BenchmarkKey -benchtime=10000x -timeout 99999s
func BenchmarkKey(b *testing.B) {
	benchTable := []struct {
		name string
	}{
		{name: "mlkem=768,mldsa=65"},
	}

	b.ResetTimer()

	for _, t := range benchTable {
		t := t
		b.Run(t.name, func(b *testing.B) {
			nowEnc := time.Now()
			for i := 0; i < b.N; i++ {
				_ = NewPrivKey()
			}
			endEnc := time.Since(nowEnc)
			b.Logf("Timer_New(N=%d): %s", b.N, endEnc)
		})
	}
}
