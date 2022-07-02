package payload

import (
	"bytes"
	"testing"
)

const (
	tcHead = 0xDEADBEAF
	tcBody = "hello, world!"
)

func TestPayload(t *testing.T) {
	pl := NewPayload(tcHead, []byte(tcBody))

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
