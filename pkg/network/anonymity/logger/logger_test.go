package logger

import (
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcService = "TST"
	tcHash    = "96cb1f0968adba001ebc216708a02c8d2817b1a77fad1206012c22716a9b130b"
	tcFmtLog  = "service=TST type=ENRSP hash=96CB1F09...6A9B130B addr=B6E23126...55F22714 proof=0000000000012345 conn=127.0.0.1:"
)

func TestLogger(t *testing.T) {
	if NewLogger("QWERTY") != nil {
		t.Error("result logger is not nil")
		return
	}

	logger := NewLogger(tcService)
	fmtLog := logger.GetFmtLog(
		CLogBaseEnqueueResponse,
		encoding.HexDecode(tcHash),
		12345,
		asymmetric.LoadRSAPubKey(testutils.TgPubKeys[0]),
		nil,
	)

	if fmtLog != tcFmtLog {
		t.Error("result fmtLog != tcFmtLog")
		return
	}
}
