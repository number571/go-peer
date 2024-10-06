package client

import (
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
BenchmarkClient/key=1024-bit-12                      100            531438 ns/op
--- BENCH: BenchmarkClient/key=1024-bit-12
    client_bench_test.go:159: Timer_Encrypt(N=1): 456.31µs
    client_bench_test.go:172: Timer_Decrypt(N=1): 450.839µs
    client_bench_test.go:159: Timer_Encrypt(N=100): 57.902233ms
    client_bench_test.go:172: Timer_Decrypt(N=100): 53.128323ms
BenchmarkClient/key=2048-bit-12                      100           1525259 ns/op
--- BENCH: BenchmarkClient/key=2048-bit-12
    client_bench_test.go:159: Timer_Encrypt(N=1): 1.40024ms
    client_bench_test.go:172: Timer_Decrypt(N=1): 1.367927ms
    client_bench_test.go:159: Timer_Encrypt(N=100): 149.796468ms
    client_bench_test.go:172: Timer_Decrypt(N=100): 152.374291ms
BenchmarkClient/key=3072-bit-12                      100           3447401 ns/op
--- BENCH: BenchmarkClient/key=3072-bit-12
    client_bench_test.go:159: Timer_Encrypt(N=1): 5.001157ms
    client_bench_test.go:172: Timer_Decrypt(N=1): 3.484949ms
    client_bench_test.go:159: Timer_Encrypt(N=100): 394.423903ms
    client_bench_test.go:172: Timer_Decrypt(N=100): 344.719275ms
BenchmarkClient/key=4096-bit-12                      100           7068881 ns/op
--- BENCH: BenchmarkClient/key=4096-bit-12
    client_bench_test.go:159: Timer_Encrypt(N=1): 7.31334ms
    client_bench_test.go:172: Timer_Decrypt(N=1): 7.212592ms
    client_bench_test.go:159: Timer_Encrypt(N=100): 805.5634ms
    client_bench_test.go:172: Timer_Decrypt(N=100): 706.869091ms
BenchmarkClient/key=5120-bit-12                      100          34466332 ns/op
--- BENCH: BenchmarkClient/key=5120-bit-12
    client_bench_test.go:159: Timer_Encrypt(N=1): 39.08925ms
    client_bench_test.go:172: Timer_Decrypt(N=1): 32.149459ms
    client_bench_test.go:159: Timer_Encrypt(N=100): 3.463671545s
    client_bench_test.go:172: Timer_Decrypt(N=100): 3.446618785s
BenchmarkClient/key=6144-bit-12                      100          58808166 ns/op
--- BENCH: BenchmarkClient/key=6144-bit-12
    client_bench_test.go:159: Timer_Encrypt(N=1): 60.5062ms
    client_bench_test.go:172: Timer_Decrypt(N=1): 62.153801ms
    client_bench_test.go:159: Timer_Encrypt(N=100): 5.954225709s
    client_bench_test.go:172: Timer_Decrypt(N=100): 5.880801097s
BenchmarkClient/key=7168-bit-12                      100          89836282 ns/op
--- BENCH: BenchmarkClient/key=7168-bit-12
    client_bench_test.go:159: Timer_Encrypt(N=1): 96.077623ms
    client_bench_test.go:172: Timer_Decrypt(N=1): 103.182745ms
    client_bench_test.go:159: Timer_Encrypt(N=100): 9.03864527s
    client_bench_test.go:172: Timer_Decrypt(N=100): 8.983613724s
BenchmarkClient/key=8192-bit-12                      100         133303055 ns/op
--- BENCH: BenchmarkClient/key=8192-bit-12
    client_bench_test.go:159: Timer_Encrypt(N=1): 135.277152ms
    client_bench_test.go:172: Timer_Decrypt(N=1): 131.425424ms
    client_bench_test.go:159: Timer_Encrypt(N=100): 13.195946659s
    client_bench_test.go:172: Timer_Decrypt(N=100): 13.330281248s
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
			name: "key=3072-bit",
			client: NewClient(
				message.NewSettings(&message.SSettings{
					FKeySizeBits:      (3 << 10),
					FMessageSizeBytes: (8 << 10),
				}),
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey3072),
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
			name: "key=5120-bit",
			client: NewClient(
				message.NewSettings(&message.SSettings{
					FKeySizeBits:      (5 << 10),
					FMessageSizeBytes: (8 << 10),
				}),
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey5120),
			),
		},
		{
			name: "key=6144-bit",
			client: NewClient(
				message.NewSettings(&message.SSettings{
					FKeySizeBits:      (6 << 10),
					FMessageSizeBytes: (8 << 10),
				}),
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey6144),
			),
		},
		{
			name: "key=7168-bit",
			client: NewClient(
				message.NewSettings(&message.SSettings{
					FKeySizeBits:      (7 << 10),
					FMessageSizeBytes: (8 << 10),
				}),
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey7168),
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
				randomBytes = append(randomBytes, testutils.PseudoRandomBytes(i))
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
