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
		t.Fatal("success create cipher with invalid key size (1)")
	}
	if cipher := NewCipherGCM([]byte{123}); cipher != nil {
		t.Fatal("success create cipher with invalid key size (2)")
	}
}

func TestEncrypt(t *testing.T) {
	t.Parallel()

	testEncrypt(t, NewCipherCFB)
	testEncrypt(t, NewCipherGCM)
}

func testEncrypt(t *testing.T, c func(pKey []byte) ICipher) {
	var (
		msg = []byte("hello, world!")
	)

	cipher := c(tgKey)

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
