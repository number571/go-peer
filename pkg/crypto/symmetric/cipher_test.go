package symmetric

import (
	"bytes"
	"testing"
)

var (
	tgKey = []byte("it is a large key with 256 bits!")
)

func TestKeySize(t *testing.T) {
	t.Parallel()

	if cipher := NewCipherCFB([]byte{123}); cipher != nil {
		t.Fatal("success create cipher with invalid key size")
	}
}

func TestEncrypt(t *testing.T) {
	t.Parallel()

	var (
		msg = []byte("hello, world!")
	)

	cipher := NewCipherCFB(tgKey)

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
