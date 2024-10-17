package client

import (
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	testutils "github.com/number571/go-peer/test/utils"
)

/*
goos: linux
goarch: amd64
pkg: github.com/number571/go-peer/pkg/client
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkClient/key=768-bit-12              1000            130936 ns/op
--- BENCH: BenchmarkClient/key=768-bit-12
    client_bench_test.go:76: Timer_Encrypt(N=1): 262.661µs
    client_bench_test.go:89: Timer_Decrypt(N=1): 177.422µs
    client_bench_test.go:76: Timer_Encrypt(N=1000): 232.488768ms
    client_bench_test.go:89: Timer_Decrypt(N=1000): 130.916195ms
PASS
*/

// go test -bench=BenchmarkClient -benchtime=100x -timeout 99999s
func BenchmarkClient(b *testing.B) {
	privKeyChain := asymmetric.NewPrivKeyChain(
		asymmetric.NewKEncPrivKey(),
		asymmetric.NewSignPrivKey(),
	)

	benchTable := []struct {
		name   string
		client IClient
	}{
		{
			name:   "kyber=768-bit,dilithium=mode3",
			client: NewClient(privKeyChain, (8 << 10)),
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
				encMsg, err := t.client.EncryptMessage(
					t.client.GetPrivKeyChain().GetKEncPrivKey().GetPubKey(),
					randomBytes[i],
				)
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
				_, _, err := t.client.DecryptMessage(encMessages[i])
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
