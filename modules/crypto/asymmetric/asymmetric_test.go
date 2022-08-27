package asymmetric

import (
	"bytes"
	"testing"
)

func TestRSASign(t *testing.T) {
	var (
		priv = NewRSAPrivKey(1024)
		msg  = []byte("hello, world!")
	)

	pub := priv.PubKey()
	sign := priv.Sign(msg)

	if !pub.Verify(msg, sign) {
		t.Error("signature is invalid")
	}
}

func TestRSAEncrypt(t *testing.T) {
	var (
		priv = NewRSAPrivKey(1024)
		msg  = []byte("hello, world!")
	)

	pub := priv.PubKey()
	emsg := pub.Encrypt(msg)

	if !bytes.Equal(msg, priv.Decrypt(emsg)) {
		t.Error("decrypted message is invalid")
	}
}
