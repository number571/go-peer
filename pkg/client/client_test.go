// nolint: goerr113
package client

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

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
