package keybuilder

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcHash = "8d47725f8604cb2e8f1be5e0b49ef143ea62625dd37ea4ce5f24501a32591784"
)

func TestKeyBuilder(t *testing.T) {
	t.Parallel()

	var (
		pasw = "hello, world!"
		salt = []byte("it's a salt!")
	)

	hash := NewKeyBuilder(testutils.TCWorkSize, salt).Build(pasw)
	if encoding.HexEncode(hash) != tcHash {
		t.Error("hash is correct?")
		return
	}

	if !bytes.Equal(hash, NewKeyBuilder(testutils.TCWorkSize, salt).Build(pasw)) {
		t.Error("hash is not determined")
		return
	}
}
