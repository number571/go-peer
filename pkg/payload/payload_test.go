package payload

import (
	"bytes"
	"testing"
)

func TestPayload64(t *testing.T) {
	t.Parallel()

	pl := NewPayload64(uint64(1), []byte("hello, world!"))

	decPl := LoadPayload64(pl.ToBytes())
	if decPl == nil {
		t.Error("decode payload is nil")
		return
	}

	if !bytes.Equal(pl.GetBody(), decPl.GetBody()) {
		t.Error("data not equal with decoded version of payload")
		return
	}

	if pl.GetHead() != decPl.GetHead() {
		t.Error("title not equal with decoded version of payload")
		return
	}

	invalidPld := LoadPayload64([]byte{1})
	if invalidPld != nil {
		t.Error("invalid payload success decoded")
		return
	}
}

func TestPayload32(t *testing.T) {
	t.Parallel()

	pl := NewPayload32(1, []byte("hello, world!"))

	decPl := LoadPayload32(pl.ToBytes())
	if decPl == nil {
		t.Error("decode payload is nil")
		return
	}

	if !bytes.Equal(pl.GetBody(), decPl.GetBody()) {
		t.Error("data not equal with decoded version of payload")
		return
	}

	if pl.GetHead() != decPl.GetHead() {
		t.Error("title not equal with decoded version of payload")
		return
	}

	invalidPld := LoadPayload32([]byte{1})
	if invalidPld != nil {
		t.Error("invalid payload success decoded")
		return
	}
}
