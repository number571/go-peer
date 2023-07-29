package logbuilder

import (
	"fmt"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcService = "TST"
	tcHash    = "96cb1f0968adba001ebc216708a02c8d2817b1a77fad1206012c22716a9b130b"
	tcFmtLog  = "service=TST type=ENRSP hash=96CB1F09...6A9B130B addr=B6E23126...55F22714 proof=0000012345 size=1024B conn=127.0.0.1:"
)

func TestLogger(t *testing.T) {
	logger := NewLogBuilder(tcService).
		WithHash(encoding.HexDecode(tcHash)).
		WithProof(12345).
		WithSize(1024).
		WithPubKey(asymmetric.LoadRSAPubKey(testutils.TgPubKeys[0]))

	if logger.Get(CLogBaseEnqueueResponse) != tcFmtLog {
		fmt.Println(logger.Get(CLogBaseEnqueueResponse))
		fmt.Println(tcFmtLog)
		t.Error("result fmtLog != tcFmtLog")
		return
	}
}
