package joiner

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
)

var tgSlice = [][]byte{
	random.NewCSPRNG().GetBytes(456),
	[]byte("hello"),
	[]byte("world->571"),
	random.NewCSPRNG().GetBytes(571),
	[]byte("qwerty"),
	{},
	random.NewCSPRNG().GetBytes(123),
	{},
}

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SJoinerError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestJoiner32(t *testing.T) {
	if _, err := LoadBytesJoiner32([]byte{1}); err == nil {
		t.Error("success load invalid bytes")
		return
	}

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
