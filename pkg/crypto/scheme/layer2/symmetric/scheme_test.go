package symmetric

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestScheme(t *testing.T) {
	t.Parallel()

	var (
		key = symmetric.NewCipherGCM(random.NewRandom().GetBytes(symmetric.CCipherKeySize))
		msg = []byte("hello, world!")
	)

	if _, err := NewScheme(8); err == nil {
		t.Fatal("success init scheme with struct size >= message size")
	}

	scheme, err := NewScheme(128)
	if err != nil {
		t.Fatal(err)
	}
	encMsg, err := scheme.EncryptMessage(key, msg)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := scheme.EncryptMessage(key, make([]byte, 256)); err == nil {
		t.Fatal("success encrypt message with overflow")
	}

	listCiphers := layer2.NewKeysContainer()
	listCiphers.Add(symmetric.NewCipherGCM(key.ToBytes()))

	gotKey, gotMsg, err := scheme.DecryptMessage(listCiphers, encMsg)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := scheme.DecryptMessage(listCiphers, []byte{}); err == nil {
		t.Fatal("success decrypt with invalid message size")
	}

	anotherListCiphers := layer2.NewKeysContainer()
	anotherListCiphers.Add(symmetric.NewCipherGCM(make([]byte, symmetric.CCipherKeySize)))
	if _, _, err := scheme.DecryptMessage(anotherListCiphers, encMsg); err == nil {
		t.Fatal("success decrypt with undefined key")
	}

	if !bytes.Equal(key.ToBytes(), gotKey.ToBytes()) {
		t.Fatal("keys are diff")
	}
	if !bytes.Equal(msg, gotMsg) {
		t.Fatal("msgs are diff")
	}

	_pubKey := asymmetric.NewPrivKey().GetPubKey()
	_keysContainer := layer2.NewKeysContainer()
	_keysContainer.Add(_pubKey)

	if _, _, err := scheme.DecryptMessage(_keysContainer, []byte{}); err == nil {
		t.Fatal("success decrypt with another key type")
	}
	if _, err := scheme.EncryptMessage(_pubKey, []byte{}); err == nil {
		t.Fatal("success encrypt with another key type")
	}

	if scheme.GetRandomKey().ToString() == scheme.GetRandomKey().ToString() { //nolint:staticcheck
		t.Fatal("random got equal values")
	}

	// fmt.Println(scheme.GetMessageSize())
	// fmt.Println(scheme.GetMessageSize() - scheme.GetPayloadLimit())
	// fmt.Println(scheme.GetPayloadLimit())
}
