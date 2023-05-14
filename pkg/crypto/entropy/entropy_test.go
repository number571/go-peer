package entropy

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestEntropy(t *testing.T) {
	var (
		msg  = []byte("hello, world!")
		salt = []byte("it's a salt!")
	)

	hash := NewEntropyBooster(testutils.TCWorkSize, salt).BoostEntropy(msg)

	if bytes.Equal(hash, hashing.NewSHA256Hasher(msg).ToBytes()) {
		t.Error("hash is correct?")
		return
	}

	if !bytes.Equal(hash, NewEntropyBooster(testutils.TCWorkSize, salt).BoostEntropy(msg)) {
		t.Error("hash is not determined")
		return
	}
}
