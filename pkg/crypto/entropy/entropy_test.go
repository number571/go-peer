package entropy

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcHash = "57daf58ec9ea7be6dbcb8a67aea76bdcfa32db40a7a3cae3d90bc977dc293599"
)

func TestEntropy(t *testing.T) {
	var (
		msg  = []byte("hello, world!")
		salt = []byte("it's a salt!")
	)

	hash := NewEntropyBooster(testutils.TCWorkSize, salt).BoostEntropy(msg)
	if encoding.HexEncode(hash) != tcHash {
		t.Error("hash is correct?")
		return
	}

	if !bytes.Equal(hash, NewEntropyBooster(testutils.TCWorkSize, salt).BoostEntropy(msg)) {
		t.Error("hash is not determined")
		return
	}
}
