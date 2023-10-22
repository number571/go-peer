package adapters

import (
	"bytes"
	"testing"
)

const (
	tcHead = 111
	tcBody = "hello, world!"
)

func TestPayload(t *testing.T) {
	t.Parallel()

	pld := NewPayload(tcHead, []byte(tcBody))

	if pld.GetHead() != tcHead {
		t.Error("pld.GetHead() != tcHead")
		return
	}

	if !bytes.Equal(pld.GetBody(), []byte(tcBody)) {
		t.Error("!bytes.Equal(pld.GetBody(), []byte(tcBody))")
		return
	}

	origPld := pld.ToOrigin()
	if origPld.GetHead() != uint64(pld.GetHead()) {
		t.Error("origPld.GetHead() != uint64(pld.GetHead())")
		return
	}
}
