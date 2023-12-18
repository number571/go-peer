package symmetric

import (
	"bytes"
	"testing"
)

var (
	tgKey = []byte("it is a large key with 256 bits!")
)

func TestAESKeySize(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	_ = NewAESCipher([]byte{123})
}

func TestAESEncrypt(t *testing.T) {
	t.Parallel()

	var (
		msg = []byte("hello, world!")
	)

	cipher := NewAESCipher(tgKey)

	emsg := cipher.EncryptBytes(msg)

	if bytes.Equal(msg, emsg) {
		t.Error("encrypted message = open message")
		return
	}

	if !bytes.Equal(msg, cipher.DecryptBytes(emsg)) {
		t.Error("decrypted message is invalid")
		return
	}

	if !bytes.Equal(cipher.DecryptBytes(emsg), cipher.DecryptBytes(emsg)) {
		t.Error("decrypted message is not determinated")
		return
	}

	if dec := cipher.DecryptBytes([]byte{123}); dec != nil {
		t.Error("success decrypt message with len < iv size")
		return
	}
}
