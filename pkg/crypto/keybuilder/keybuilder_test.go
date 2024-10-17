package keybuilder

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
)

const (
	tcKeySize = 32
	tcHash    = "0aac5a0a082dcbc64befb61c798549126cd6b4075081debc6636c53c941c5cb0"
)

func TestKeyBuilder(t *testing.T) {
	t.Parallel()

	var (
		pasw = "hello, world!"
		salt = []byte("it's a salt!")
	)

	hash := NewKeyBuilder(1<<10, salt).Build(pasw, tcKeySize)
	if encoding.HexEncode(hash) != tcHash {
		t.Error("hash is correct?")
		return
	}

	if !bytes.Equal(hash, NewKeyBuilder(1<<10, salt).Build(pasw, tcKeySize)) {
		t.Error("hash is not determined")
		return
	}
}
