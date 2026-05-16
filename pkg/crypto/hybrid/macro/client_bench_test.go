package macro

import (
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hybrid"
	testutils "github.com/number571/go-peer/test/utils"
)

/*
goos: linux
goarch: amd64
pkg: github.com/number571/go-peer/pkg/scheme
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkScheme/mlkem=768,mldsa=65-12              10000             94239 ns/op
--- BENCH: BenchmarkScheme/mlkem=768,mldsa=65-12
    scheme_bench_test.go:67: Timer_Encrypt(N=1): 207.249µs
    scheme_bench_test.go:80: Timer_Decrypt(N=1): 131.503µs
    scheme_bench_test.go:67: Timer_Encrypt(N=10000): 2.376859383s
    scheme_bench_test.go:80: Timer_Decrypt(N=10000): 942.366433ms
PASS
*/

// go test -bench=BenchmarkScheme -benchtime=1000x -timeout 99999s
func BenchmarkScheme(b *testing.B) {
	privKey := asymmetric.NewPrivKey()
	pubKey := privKey.GetPubKey()

	mapKeys := asymmetric.NewMapPubKeys()
	mapKeys.SetPubKey(pubKey)

	benchTable := []struct {
		name   string
		scheme hybrid.IScheme
	}{
		{
			name:   "mlkem=768,mldsa=65",
			scheme: NewScheme(privKey, (8 << 10)),
		},
	}

	b.ResetTimer()

	for _, t := range benchTable {
		t := t
		b.Run(t.name, func(b *testing.B) {
			b.StopTimer()
			encMessages := make([][]byte, b.N)
			randomBytes := make([][]byte, 0, b.N)
			for i := 0; i < b.N; i++ {
				randomBytes = append(randomBytes, testutils.PseudoRandomBytes(i))
			}
			b.StartTimer()

			nowEnc := time.Now()
			for i := 0; i < b.N; i++ {
				encMsg, err := t.scheme.EncryptMessage(pubKey, randomBytes[i])
				if err != nil {
					b.Error(err)
					return
				}
				encMessages[i] = encMsg
			}
			endEnc := time.Since(nowEnc)

			b.Logf("Timer_Encrypt(N=%d): %s", b.N, endEnc)
			b.ResetTimer()

			nowDec := time.Now()
			for i := 0; i < b.N; i++ {
				_, _, err := t.scheme.DecryptMessage(mapKeys, encMessages[i])
				if err != nil {
					b.Error(err)
					return
				}
			}
			endDec := time.Since(nowDec)

			b.Logf("Timer_Decrypt(N=%d): %s", b.N, endDec)
		})
	}
}
