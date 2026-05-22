package symmetric

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestPanicNewScheme(t *testing.T) {
	t.Parallel()

	tcNewSchemeWithSmallMsgSize(t)
}

func tcNewSchemeWithSmallMsgSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = NewScheme(8)
}

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
	if _, err := scheme.EncryptMessage(key, make([]byte, 256)); err == nil {
		t.Fatal("success encrypt message with overflow")
	}

	listCiphers := symmetric.NewListCiphers()
	listCiphers.Add(symmetric.NewCipherGCM(key.ToBytes()))

	gotKey, gotMsg, err := scheme.DecryptMessage(listCiphers, encMsg)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := scheme.DecryptMessage(listCiphers, []byte{}); err == nil {
		t.Fatal("success decrypt with invalid message size")
	}

	anotherListCiphers := symmetric.NewListCiphers()
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

	if _, _, err := scheme.DecryptMessage(asymmetric.NewMapPubKeys(), []byte{}); err == nil {
		t.Error("success decrypt with another key type")
		return
	}
	if _, err := scheme.EncryptMessage(asymmetric.NewPrivKey().GetPubKey(), []byte{}); err == nil {
		t.Error("success encrypt with another key type")
		return
	}

	if scheme.GetRandomKey().ToString() == scheme.GetRandomKey().ToString() { //nolint:staticcheck
		t.Fatal("random got equal values")
	}

	// fmt.Println(scheme.GetMessageSize())
	// fmt.Println(scheme.GetMessageSize() - scheme.GetPayloadLimit())
	// fmt.Println(scheme.GetPayloadLimit())
}
