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
		t.Fatal("decode payload is nil")
	}

	if !bytes.Equal(pl.GetBody(), decPl.GetBody()) {
		t.Fatal("data not equal with decoded version of payload")
	}

	if pl.GetHead() != decPl.GetHead() {
		t.Fatal("title not equal with decoded version of payload")
	}

	invalidPld := LoadPayload64([]byte{1})
	if invalidPld != nil {
		t.Fatal("invalid payload success decoded")
	}
}

func TestPayload32(t *testing.T) {
	t.Parallel()

	pl := NewPayload32(1, []byte("hello, world!"))

	decPl := LoadPayload32(pl.ToBytes())
	if decPl == nil {
		t.Fatal("decode payload is nil")
	}

	if !bytes.Equal(pl.GetBody(), decPl.GetBody()) {
		t.Fatal("data not equal with decoded version of payload")
	}

	if pl.GetHead() != decPl.GetHead() {
		t.Fatal("title not equal with decoded version of payload")
	}

	invalidPld := LoadPayload32([]byte{1})
	if invalidPld != nil {
		t.Fatal("invalid payload success decoded")
	}
}
