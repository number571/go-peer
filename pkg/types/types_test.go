package types

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
)

type tsSomeStruct struct {
	V string `json:"v"`
	N int    `json:"n"`
}

func TestPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()

	_ = NewConverter(func() {})
}

func TestConverter(t *testing.T) {
	t.Parallel()

	b := []byte{111, 222, 123}
	conv1 := NewConverter(b)
	if !bytes.Equal(b, conv1.ToBytes()) {
		t.Fatal("invalid bytes (1)")
	}
	if conv1.ToString() != encoding.HexEncode(b) {
		t.Fatal("invalid string (1)")
	}

	s := "hello, world!"
	conv2 := NewConverter(s)
	if !bytes.Equal([]byte(s), conv2.ToBytes()) {
		t.Fatal("invalid bytes (2)")
	}
	if conv2.ToString() != s {
		t.Fatal("invalid string (2)")
	}

	v := &tsSomeStruct{V: "abc", N: 123}
	conv3 := NewConverter(v)
	if !bytes.Equal([]byte(`{"v":"abc","n":123}`), conv3.ToBytes()) {
		t.Fatal("invalid bytes (3)")
	}
	if conv3.ToString() != `{"v":"abc","n":123}` {
		t.Fatal("invalid string (3)")
	}
}
