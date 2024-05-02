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
	{},
	random.NewStdPRNG().GetBytes(123),
	{},
}

func TestError(t *testing.T) {
	str := "value"
	err := &SJoinerError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestJoiner32(t *testing.T) {
	joinerBytes := NewBytesJoiner32(tgSlice)

	slice, err := LoadBytesJoiner32(joinerBytes)
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
