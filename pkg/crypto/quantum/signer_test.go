package quantum

import (
	"testing"
)

func TestSigner(t *testing.T) {
	t.Parallel()

	privKey := NewSignerPrivKey()
	privKey = LoadSignerPrivKey(privKey.ToBytes())

	pubKey := privKey.GetPubKey()
	pubKey = LoadSignerPubKey(pubKey.ToBytes())

	msg := []byte("hello, world!")
	sign := privKey.SignBytes(msg)

	if !pubKey.VerifyBytes(msg, sign) {
		t.Error("invalid verify")
		return
	}

	// fmt.Println(len(pubKey.ToBytes()), len(sign))
}
