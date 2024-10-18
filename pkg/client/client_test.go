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

	kemPubKey := client.GetPrivKey().GetKEncPrivKey().GetPubKey()
	msg := []byte("hello, world!")

	enc, err := client.EncryptMessage(kemPubKey, msg)
	if err != nil {
		t.Error(err)
		return
	}

	// os.WriteFile("test_binary.msg", enc, 0600)
	// os.WriteFile("test_string.msg", []byte(encoding.HexEncode(enc)), 0600)

	signerPubKey := client.GetPrivKey().GetSignPrivKey().GetPubKey()
	gotSignPubKey, dec, err := client.DecryptMessage(enc)
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(signerPubKey.ToBytes(), gotSignPubKey.ToBytes()) {
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
