package joiner

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
)

var tgSlice = [][]byte{
	random.NewRandom().GetBytes(456),
	[]byte("hello"),
	[]byte("world->571"),
	random.NewRandom().GetBytes(571),
	[]byte("qwerty"),
	{},
	random.NewRandom().GetBytes(123),
	{},
}

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SJoinerError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestJoiner32(t *testing.T) {
	if _, err := LoadBytesJoiner32([]byte{1}); err == nil {
		t.Fatal("success load invalid bytes")
	}

	joinerBytes := NewBytesJoiner32(tgSlice)

	slice, err := LoadBytesJoiner32(joinerBytes)
	if err != nil {
		t.Fatal(err)
	}

	if len(slice) != len(tgSlice) {
		t.Fatal("len(slice) != len(tgSlice)")
	}

	for i := range slice {
		if !bytes.Equal(slice[i], tgSlice[i]) {
			t.Fatal("!bytes.Equal(slice[i],tgSlice[i])")
		}
	}
}
