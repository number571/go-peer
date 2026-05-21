package symmetric

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
)

func TestScheme(t *testing.T) {
	t.Parallel()

	var (
		key = symmetric.NewCipherGCM(random.NewRandom().GetBytes(symmetric.CCipherKeySize))
		msg = []byte("hello, world!")
	)

	scheme := NewScheme(128)
	encMsg, err := scheme.EncryptMessage(key, msg)
	if err != nil {
		t.Fatal(err)
	}

	listCiphers := symmetric.NewListCiphers()
	listCiphers.Add(symmetric.NewCipherGCM(key.ToBytes()))

	gotKey, gotMsg, err := scheme.DecryptMessage(listCiphers, encMsg)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(key.ToBytes(), gotKey.ToBytes()) {
		t.Fatal("keys are diff")
	}
	if !bytes.Equal(msg, gotMsg) {
		fmt.Println(string(gotMsg))
		t.Fatal("msgs are diff")
	}

	fmt.Println(scheme.GetMessageSize())
	fmt.Println(scheme.GetMessageSize() - scheme.GetPayloadLimit())
	fmt.Println(scheme.GetPayloadLimit())
}
