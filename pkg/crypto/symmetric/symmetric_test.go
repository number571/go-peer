package symmetric

import (
	"bytes"
	"testing"
)

func TestAESEncrypt(t *testing.T) {
	var (
		key = []byte("it's a key!")
		msg = []byte("hello, world!")
	)

	cipher := NewAESCipher(key)
	emsg := cipher.EncryptBytes(msg)

	if bytes.Equal(msg, emsg) {
		t.Error("encrypted message = open message")
	}

	if !bytes.Equal(msg, cipher.DecryptBytes(emsg)) {
		t.Error("decrypted message is invalid")
	}

	if !bytes.Equal(cipher.DecryptBytes(emsg), cipher.DecryptBytes(emsg)) {
		t.Error("decrypted message is not determinated")
	}
}
