package payload

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/utils/testutils"
)

func TestPayload(t *testing.T) {
	pl := NewPayload(uint64(testutils.TcHead), []byte(testutils.TcBody))

	decPl := LoadPayload(pl.Bytes())
	if decPl == nil {
		t.Error("decode payload is nil")
		return
	}

	if !bytes.Equal(pl.Body(), decPl.Body()) {
		t.Error("data not equal with decoded version of payload")
		return
	}

	if pl.Head() != decPl.Head() {
		t.Error("title not equal with decoded version of payload")
		return
	}
}
