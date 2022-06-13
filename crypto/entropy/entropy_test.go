package entropy

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/crypto/hashing"
)

func TestEntropy(t *testing.T) {
	var (
		msg  = []byte("hello, world!")
		salt = []byte("it's a salt!")
	)

	hash := NewEntropy(10).Raise(msg, salt)

	if bytes.Equal(hash, hashing.NewSHA256Hasher(msg).Bytes()) {
		t.Error("hash is correct?")
	}

	if !bytes.Equal(hash, NewEntropy(10).Raise(msg, salt)) {
		t.Error("hash is not determined")
	}
}
