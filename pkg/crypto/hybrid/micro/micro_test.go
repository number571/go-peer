package micro

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/hybrid"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/types"
)

func TestScheme(t *testing.T) {
	t.Parallel()

	var (
		key = types.NewConverter(random.NewRandom().GetBytes(symmetric.CCipherKeySize))
		msg = []byte("hello, world!")
	)

	scheme := NewScheme(200)
	encMsg, err := scheme.EncryptMessage(key, msg)
	if err != nil {
		t.Fatal(err)
	}

	gotKey, gotMsg, err := scheme.DecryptMessage([]hybrid.IParticipantKey{key}, encMsg)
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

	// fmt.Println(scheme.GetMessageSize())
	// fmt.Println(scheme.GetMessageSize() - scheme.GetPayloadLimit())
	// fmt.Println(scheme.GetPayloadLimit())
}
