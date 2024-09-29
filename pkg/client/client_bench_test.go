package client

import (
	"math/rand"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	testutils "github.com/number571/go-peer/test/utils"
)

/*
goos: linux
goarch: amd64
pkg: github.com/number571/go-peer/pkg/client
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkClient/key=1024-bit-12               100            433814 ns/op
--- BENCH: BenchmarkClient/key=1024-bit-12
    client_bench_test.go:96: Timer_Encrypt(N=1): 458.299µs
    client_bench_test.go:109: Timer_Decrypt(N=1): 408.396µs
    client_bench_test.go:96: Timer_Encrypt(N=100): 60.532764ms
    client_bench_test.go:109: Timer_Decrypt(N=100): 43.366633ms
BenchmarkClient/key=2048-bit-12               100           1296076 ns/op
--- BENCH: BenchmarkClient/key=2048-bit-12
    client_bench_test.go:96: Timer_Encrypt(N=1): 1.328712ms
    client_bench_test.go:109: Timer_Decrypt(N=1): 1.243297ms
    client_bench_test.go:96: Timer_Encrypt(N=100): 134.830767ms
    client_bench_test.go:109: Timer_Decrypt(N=100): 129.589976ms
BenchmarkClient/key=4096-bit-12               100           6924869 ns/op
--- BENCH: BenchmarkClient/key=4096-bit-12
    client_bench_test.go:96: Timer_Encrypt(N=1): 7.454978ms
    client_bench_test.go:109: Timer_Decrypt(N=1): 8.154685ms
    client_bench_test.go:96: Timer_Encrypt(N=100): 726.880386ms
    client_bench_test.go:109: Timer_Decrypt(N=100): 692.402735ms
BenchmarkClient/key=8192-bit-12               100         129773452 ns/op
--- BENCH: BenchmarkClient/key=8192-bit-12
    client_bench_test.go:96: Timer_Encrypt(N=1): 128.75005ms
    client_bench_test.go:109: Timer_Decrypt(N=1): 128.397816ms
    client_bench_test.go:96: Timer_Encrypt(N=100): 12.717255259s
    client_bench_test.go:109: Timer_Decrypt(N=100): 12.97728186s
PASS
*/

// go test -bench=BenchmarkClient -benchtime=100x -timeout 99999s
func BenchmarkClient(b *testing.B) {
	benchTable := []struct {
		name   string
		client IClient
	}{
		{
			name: "key=1024-bit",
			client: NewClient(
				message.NewSettings(&message.SSettings{
					FKeySizeBits:      (1 << 10),
					FMessageSizeBytes: (8 << 10),
				}),
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey1024),
			),
		},
		{
			name: "key=2048-bit",
			client: NewClient(
				message.NewSettings(&message.SSettings{
					FKeySizeBits:      (2 << 10),
					FMessageSizeBytes: (8 << 10),
				}),
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey2048),
			),
		},
		{
			name: "key=4096-bit",
			client: NewClient(
				message.NewSettings(&message.SSettings{
					FKeySizeBits:      (4 << 10),
					FMessageSizeBytes: (8 << 10),
				}),
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey4096),
			),
		},
		{
			name: "key=8192-bit",
			client: NewClient(
				message.NewSettings(&message.SSettings{
					FKeySizeBits:      (8 << 10),
					FMessageSizeBytes: (8 << 10),
				}),
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey8192),
			),
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
				randomBytes = append(randomBytes, testPseudoRandomBytes(i))
			}
			b.StartTimer()

			nowEnc := time.Now()
			for i := 0; i < b.N; i++ {
				encMsg, err := t.client.EncryptMessage(
					t.client.GetPubKey(),
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

func testPseudoRandomBytes(pSeed int) []byte {
	r := rand.New(rand.NewSource(int64(pSeed))) //nolint:gosec
	result := make([]byte, 0, 16)
	for i := 0; i < 16; i++ {
		result = append(result, byte(r.Intn(256)))
	}
	return result
}
