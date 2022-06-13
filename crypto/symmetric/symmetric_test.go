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
	emsg := cipher.Encrypt(msg)

	if !bytes.Equal(msg, cipher.Decrypt(emsg)) {
		t.Error("decrypted message is invalid")
	}
}
