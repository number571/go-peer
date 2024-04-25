package joiner

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
)

var tgSlice = [][]byte{
	random.NewStdPRNG().GetBytes(456),
	[]byte("hello"),
	[]byte("world->571"),
	random.NewStdPRNG().GetBytes(571),
	[]byte("qwerty"),
	random.NewStdPRNG().GetBytes(123),
}

func TestJoiner(t *testing.T) {
	joinerBytes := NewBytesJoiner(tgSlice)

	slice, err := LoadBytesJoiner(joinerBytes)
	if err != nil {
		t.Error(err)
		return
	}

	if len(slice) != len(tgSlice) {
		t.Error("len(slice) != len(tgSlice)")
		return
	}

	for i := range slice {
		if !bytes.Equal(slice[i], tgSlice[i]) {
			t.Error("!bytes.Equal(slice[i],tgSlice[i])")
			return
		}
	}
}
