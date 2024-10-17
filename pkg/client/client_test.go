// nolint: goerr113
package client

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func TestClient(t *testing.T) {
	t.Parallel()

	client := NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: (8 << 10),
			FEncKeySizeBytes:  asymmetric.CKEncSize,
		}),
		asymmetric.NewPrivKeyChain(
			asymmetric.NewKEncPrivKey(),
			asymmetric.NewSignPrivKey(),
		),
	)

	kemPubKey := client.GetPrivKeyChain().GetKEncPrivKey().GetPubKey()
	msg := []byte("hello, world!")

	enc, err := client.EncryptMessage(kemPubKey, msg)
	if err != nil {
		t.Error(err)
		return
	}

	signerPubKey := client.GetPrivKeyChain().GetSignPrivKey().GetPubKey()
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

	// fmt.Println(client.GetMessageLimit())
	// fmt.Println(len(client.GetPrivKeyChain().GetPubKeyChain().ToString()))
}
