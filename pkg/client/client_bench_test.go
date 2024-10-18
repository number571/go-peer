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
BenchmarkClient/kyber=768-bit,dilithium=mode3-12                    1000            132328 ns/op
--- BENCH: BenchmarkClient/kyber=768-bit,dilithium=mode3-12
    client_bench_test.go:66: Timer_Encrypt(N=1): 280.196µs
    client_bench_test.go:79: Timer_Decrypt(N=1): 155.822µs
    client_bench_test.go:66: Timer_Encrypt(N=1000): 254.339142ms
    client_bench_test.go:79: Timer_Decrypt(N=1000): 132.305197ms
PASS
*/

// go test -bench=BenchmarkClient -benchtime=1000x -timeout 99999s
func BenchmarkClient(b *testing.B) {
	privKeyChain := asymmetric.NewPrivKey()

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
					t.client.GetPrivKey().GetKEncPrivKey().GetPubKey(),
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
