// nolint: goerr113
package client

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func TestPanicNewClient(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = NewClient(asymmetric.NewPrivKey(), 8)
}

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SClientError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestClient(t *testing.T) {
	t.Parallel()

	client := NewClient(
		asymmetric.NewPrivKey(),
		(8 << 10),
	)

	kemPubKey := client.GetPrivKey().GetKEMPrivKey().GetPubKey()
	msg := []byte("hello, world!")

	enc, err := client.EncryptMessage(kemPubKey, msg)
	if err != nil {
		t.Error(err)
		return
	}

	// _ = os.WriteFile("message/test_binary.msg", enc, 0600)
	// _ = os.WriteFile("message/test_string.msg", []byte(encoding.HexEncode(enc)), 0600)

	signerPubKey := client.GetPrivKey().GetDSAPrivKey().GetPubKey()
	gotDSAPubKey, dec, err := client.DecryptMessage(enc)
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(signerPubKey.ToBytes(), gotDSAPubKey.ToBytes()) {
		t.Error("invalid decrypt signer key")
		return
	}
	if !bytes.Equal(msg, dec) {
		t.Error("invalid decrypt message")
		return
	}

	// fmt.Println(client.GetPayloadLimit(), client.GetMessageSize())
	// fmt.Println(len(client.GetPrivKey().GetPubKey().ToString()))
}

func TestDecrypt(t *testing.T) {
	t.Parallel()

	client := NewClient(
		asymmetric.NewPrivKey(),
		(8 << 10),
	)

	if _, _, err := client.DecryptMessage([]byte{123}); err == nil {
		t.Error("success decrypt with invalid ciphertext (1)")
		return
	}

	kemPubKey := client.GetPrivKey().GetKEMPrivKey().GetPubKey()
	msg := []byte("hello, world!")

	enc, err := client.EncryptMessage(kemPubKey, msg)
	if err != nil {
		t.Error(err)
		return
	}

	enc[0] ^= 1
	if _, _, err := client.DecryptMessage(enc); err == nil {
		t.Error("success decrypt with invalid ciphertext (2)")
		return
	}

	enc[0] ^= 1
	enc[len(enc)-1] ^= 1
	if _, _, err := client.DecryptMessage(enc); err == nil {
		t.Error("success decrypt with invalid ciphertext (3)")
		return
	}

	enc[len(enc)-1] ^= 1
	enc[len(enc)-2000] ^= 1
	if _, _, err := client.DecryptMessage(enc); err == nil {
		t.Error("success decrypt with invalid ciphertext (4)")
		return
	}
}
