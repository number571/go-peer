package layer1

import (
	"runtime"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/utils"
)

/*
goos: linux
goarch: amd64
pkg: github.com/number571/go-peer/pkg/network/message
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkMessageSequence/worksize=20-bit-12                  100             29079 ns/op
--- BENCH: BenchmarkMessageSequence/worksize=20-bit-12
    message_bench_test.go:66: Timer_New(N=1): 587.357878ms
    message_bench_test.go:79: Timer_Load(N=1): 26.216µs
    message_bench_test.go:66: Timer_New(N=100): 47.755895776s
    message_bench_test.go:79: Timer_Load(N=100): 2.895318ms
BenchmarkMessageSequence/worksize=21-bit-12                  100             17722 ns/op
--- BENCH: BenchmarkMessageSequence/worksize=21-bit-12
    message_bench_test.go:66: Timer_New(N=1): 862.067758ms
    message_bench_test.go:79: Timer_Load(N=1): 31.847µs
    message_bench_test.go:66: Timer_New(N=100): 1m32.635054736s
    message_bench_test.go:79: Timer_Load(N=100): 1.762232ms
BenchmarkMessageSequence/worksize=22-bit-12                  100             24118 ns/op
--- BENCH: BenchmarkMessageSequence/worksize=22-bit-12
    message_bench_test.go:66: Timer_New(N=1): 687.835779ms
    message_bench_test.go:79: Timer_Load(N=1): 29.338µs
    message_bench_test.go:66: Timer_New(N=100): 3m18.467072934s
    message_bench_test.go:79: Timer_Load(N=100): 2.399807ms
BenchmarkMessageSequence/worksize=23-bit-12                  100             18015 ns/op
--- BENCH: BenchmarkMessageSequence/worksize=23-bit-12
    message_bench_test.go:66: Timer_New(N=1): 825.817871ms
    message_bench_test.go:79: Timer_Load(N=1): 86.881µs
    message_bench_test.go:66: Timer_New(N=100): 6m0.166963461s
    message_bench_test.go:79: Timer_Load(N=100): 1.791463ms
PASS
*/

// go test -bench=BenchmarkMessageSequence -benchtime=100x -timeout 99999s
func BenchmarkMessageSequence(b *testing.B) {
	benchTable := []struct {
		name string
		sett IConstructSettings
	}{
		{
			name: "worksize=20-bit",
			sett: NewConstructSettings(&SConstructSettings{
				FSettings: NewSettings(&SSettings{
					FWorkSizeBits: 20,
				}),
			}),
		},
		{
			name: "worksize=21-bit",
			sett: NewConstructSettings(&SConstructSettings{
				FSettings: NewSettings(&SSettings{
					FWorkSizeBits: 21,
				}),
			}),
		},
		{
			name: "worksize=22-bit",
			sett: NewConstructSettings(&SConstructSettings{
				FSettings: NewSettings(&SSettings{
					FWorkSizeBits: 22,
				}),
			}),
		},
		{
			name: "worksize=23-bit",
			sett: NewConstructSettings(&SConstructSettings{
				FSettings: NewSettings(&SSettings{
					FWorkSizeBits: 23,
				}),
			}),
		},
	}

	b.ResetTimer()

	for _, t := range benchTable {
		t := t
		b.Run(t.name, func(b *testing.B) {
			b.StopTimer()
			messages := make([]IMessage, b.N)
			randomPayloads := make([]payload.IPayload32, 0, b.N)
			for i := 0; i < b.N; i++ {
				randomPayloads = append(
					randomPayloads,
					payload.NewPayload32(1, testutils.PseudoRandomBytes(i)),
				)
			}
			b.StartTimer()

			nowNew := time.Now()
			for i := 0; i < b.N; i++ {
				messages[i] = NewMessage(t.sett, randomPayloads[i])
			}
			endNew := time.Since(nowNew)

			b.Logf("Timer_New(N=%d): %s", b.N, endNew)
			b.ResetTimer()

			nowLoad := time.Now()
			for i := 0; i < b.N; i++ {
				_, err := LoadMessage(t.sett.GetSettings(), messages[i].ToBytes())
				if err != nil {
					b.Error(err)
					return
				}
			}
			endLoad := time.Since(nowLoad)

			b.Logf("Timer_Load(N=%d): %s", b.N, endLoad)
		})
	}
}

/*
goos: linux
goarch: amd64
pkg: github.com/number571/go-peer/pkg/network/message
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkMessageParallel/worksize=20-bit-12                100             29409 ns/op
--- BENCH: BenchmarkMessageParallel/worksize=20-bit-12
    message_bench_test.go:130: Timer_New(N=1): 150.160614ms
    message_bench_test.go:143: Timer_Load(N=1): 106.167µs
    message_bench_test.go:130: Timer_New(N=100): 13.159000387s
    message_bench_test.go:143: Timer_Load(N=100): 2.928936ms
BenchmarkMessageParallel/worksize=21-bit-12                100             34238 ns/op
--- BENCH: BenchmarkMessageParallel/worksize=21-bit-12
    message_bench_test.go:130: Timer_New(N=1): 199.56114ms
    message_bench_test.go:143: Timer_Load(N=1): 74.731µs
    message_bench_test.go:130: Timer_New(N=100): 25.794968877s
    message_bench_test.go:143: Timer_Load(N=100): 3.411114ms
BenchmarkMessageParallel/worksize=22-bit-12                100             32504 ns/op
--- BENCH: BenchmarkMessageParallel/worksize=22-bit-12
    message_bench_test.go:130: Timer_New(N=1): 227.384182ms
    message_bench_test.go:143: Timer_Load(N=1): 19.048µs
    message_bench_test.go:130: Timer_New(N=100): 58.771191158s
    message_bench_test.go:143: Timer_Load(N=100): 3.23526ms
BenchmarkMessageParallel/worksize=23-bit-12                100             32992 ns/op
--- BENCH: BenchmarkMessageParallel/worksize=23-bit-12
    message_bench_test.go:130: Timer_New(N=1): 245.354315ms
    message_bench_test.go:143: Timer_Load(N=1): 51.538µs
    message_bench_test.go:130: Timer_New(N=100): 1m54.765028509s
    message_bench_test.go:143: Timer_Load(N=100): 3.282597ms
PASS
*/

// go test -bench=BenchmarkMessageParallel -benchtime=100x -timeout 99999s
func BenchmarkMessageParallel(b *testing.B) {
	gomaxprocs := uint64(runtime.GOMAXPROCS(0)) //nolint:gosec

	benchTable := []struct {
		name string
		sett IConstructSettings
	}{
		{
			name: "worksize=20-bit",
			sett: NewConstructSettings(&SConstructSettings{
				FSettings: NewSettings(&SSettings{
					FWorkSizeBits: 20,
				}),
				FParallel: gomaxprocs,
			}),
		},
		{
			name: "worksize=21-bit",
			sett: NewConstructSettings(&SConstructSettings{
				FSettings: NewSettings(&SSettings{
					FWorkSizeBits: 21,
				}),
				FParallel: gomaxprocs,
			}),
		},
		{
			name: "worksize=22-bit",
			sett: NewConstructSettings(&SConstructSettings{
				FSettings: NewSettings(&SSettings{
					FWorkSizeBits: 22,
				}),
				FParallel: gomaxprocs,
			}),
		},
		{
			name: "worksize=23-bit",
			sett: NewConstructSettings(&SConstructSettings{
				FSettings: NewSettings(&SSettings{
					FWorkSizeBits: 23,
				}),
				FParallel: gomaxprocs,
			}),
		},
	}

	b.ResetTimer()

	for _, t := range benchTable {
		t := t
		b.Run(t.name, func(b *testing.B) {
			b.StopTimer()
			messages := make([]IMessage, b.N)
			randomPayloads := make([]payload.IPayload32, 0, b.N)
			for i := 0; i < b.N; i++ {
				randomPayloads = append(
					randomPayloads,
					payload.NewPayload32(1, testutils.PseudoRandomBytes(i)),
				)
			}
			b.StartTimer()

			nowNew := time.Now()
			for i := 0; i < b.N; i++ {
				messages[i] = NewMessage(t.sett, randomPayloads[i])
			}
			endNew := time.Since(nowNew)

			b.Logf("Timer_New(N=%d): %s", b.N, endNew)
			b.ResetTimer()

			nowLoad := time.Now()
			for i := 0; i < b.N; i++ {
				_, err := LoadMessage(t.sett.GetSettings(), messages[i].ToBytes())
				if err != nil {
					b.Error(err)
					return
				}
			}
			endLoad := time.Since(nowLoad)

			b.Logf("Timer_Load(N=%d): %s", b.N, endLoad)
		})
	}
}
