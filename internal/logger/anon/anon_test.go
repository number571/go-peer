package anon

import (
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcService = "TST"
	tcHash    = "96cb1f0968adba001ebc216708a02c8d2817b1a77fad1206012c22716a9b130b"
	tcFmtLog  = "service=TST type=ENQRQ hash=96CB1F09...6A9B130B addr=6245E00D...327047E9 proof=0000012345 size=1024B conn=127.0.0.1:"
)

func TestLoggerPanic(t *testing.T) {
	t.Parallel()

	logFunc := GetLogFunc()
	for i := 0; i < 3; i++ {
		testLoggerPanic(t, logFunc, i)
	}
}

func testLoggerPanic(t *testing.T, f logger.ILogFunc, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	switch n {
	case 0:
		f(struct{}{})
	case 1:
		logger := testNewAnonLogger()
		f(logger) // without type
	case 2:
		logger := testNewAnonLogger().WithType(255)
		f(logger) // with unknown type
	}
}

func TestLogger(t *testing.T) {
	t.Parallel()

	logger := testNewAnonLogger().
		WithType(anon_logger.CLogBaseEnqueueRequest)

	logFunc := GetLogFunc()
	if logFunc(logger) != tcFmtLog {
		t.Error("result fmtLog != tcFmtLog")
		return
	}
}

func testNewAnonLogger() anon_logger.ILogBuilder {
	return anon_logger.NewLogBuilder(tcService).
		WithHash(encoding.HexDecode(tcHash)).
		WithProof(12345).
		WithSize(1024).
		WithPubKey(asymmetric.LoadRSAPubKey(testutils.TgPubKeys[0]))
}
