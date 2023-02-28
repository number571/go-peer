package entropy

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/hashing"
)

func TestEntropy(t *testing.T) {
	var (
		msg  = []byte("hello, world!")
		salt = []byte("it's a salt!")
	)

	hash := NewEntropyBooster(10, salt).BoostEntropy(msg)

	if bytes.Equal(hash, hashing.NewSHA256Hasher(msg).ToBytes()) {
		t.Error("hash is correct?")
	}

	if !bytes.Equal(hash, NewEntropyBooster(10, salt).BoostEntropy(msg)) {
		t.Error("hash is not determined")
	}
}
