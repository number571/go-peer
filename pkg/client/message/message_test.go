package message

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcSalt    = "327d3d5fa3e8df5012ca12a8d65d0644669861797f659ac13adf15c2508bec5bd43538ba6c5c6871992c35ebafdaf019"
	tcSession = "58e578980580a793c01765b6d8709784d41dbcb53f3fe87e7e4e4d4b33c4f93295bac351966d6cd70f5e352cae1d290bec3b8234d4048807e5aa4d3bb983ab0d668be1f2843b77aea181ea8b92ffbf641aa5f8e25e12243ade696cbace216c5148cb68d0ae93342da9712561b1ca806dd992ce6d14b1cb0d6211780faae7ca15"
	tcSender  = "ac2a318e0e0f4548f07b5444d293304b52bde3e22e23da82147d03a4b6721b98dc43a53adbc1ff33a3819f0c3899b3e2dc9b7cb18b0da96ad6c4ff93123764b68ec660314ee13167c7ece3bb73c72258e05c273155b2bb13860820aebf9c9400158f38b3208826e176ef5b90ea209c5fc302e8162f75ae8fa2bf4c50f5c16256d5d413e6bf69592ed90ae9f91af590c06faf71b5e87ba3dc4a4f96b3"

	tcSign = "cd18d99ec9af3dabd548a5e2b8e4f34add2605671c733342968159eb9fc23484ca14b7b6e465fa3527bcc03a50aa6203240077c5118e18f28d48939df5bfab4fa3d59ed3b23a0bc6a8ba08275256a6c3577bcb887ed7657578e7de711eac7af88de9f75f0fce0428134aa24fe5c9655256c7753d54b49a5504c8d1a2a1d242c8d4459431474ebe9f5e4bd005d7a94681"
	tcHash = "3d90d389b7454bd3f1a84c4509c0afdc2a92d6acafdd387a541590ed6caf34e0"

	// -1 byte from hash
	tcInvalidPrefix = `{"pubk":"ac2a318e0e0f4548f07b5444d293304b52bde3e22e23da82147d03a4b6721b98dc43a53adbc1ff33a3819f0c3899b3e2dc9b7cb18b0da96ad6c4ff93123764b68ec660314ee13167c7ece3bb73c72258e05c273155b2bb13860820aebf9c9400158f38b3208826e176ef5b90ea209c5fc302e8162f75ae8fa2bf4c50f5c16256d5d413e6bf69592ed90ae9f91af590c06faf71b5e87ba3dc4a4f96b3","enck":"58e578980580a793c01765b6d8709784d41dbcb53f3fe87e7e4e4d4b33c4f93295bac351966d6cd70f5e352cae1d290bec3b8234d4048807e5aa4d3bb983ab0d668be1f2843b77aea181ea8b92ffbf641aa5f8e25e12243ade696cbace216c5148cb68d0ae93342da9712561b1ca806dd992ce6d14b1cb0d6211780faae7ca15","salt":"327d3d5fa3e8df5012ca12a8d65d0644669861797f659ac13adf15c2508bec5bd43538ba6c5c6871992c35ebafdaf019","hash":"3dd389b7454bd3f1a84c4509c0afdc2a92d6acafdd387a541590ed6caf34e0","sign":"cd18d99ec9af3dabd548a5e2b8e4f34add2605671c733342968159eb9fc23484ca14b7b6e465fa3527bcc03a50aa6203240077c5118e18f28d48939df5bfab4fa3d59ed3b23a0bc6a8ba08275256a6c3577bcb887ed7657578e7de711eac7af88de9f75f0fce0428134aa24fe5c9655256c7753d54b49a5504c8d1a2a1d242c8d4459431474ebe9f5e4bd005d7a94681"}@`
)

func TestSettings(t *testing.T) {
	t.Parallel()

	for i := 0; i < 1; i++ {
		testSettings(t, i)
	}
}

func testSettings(t *testing.T, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	switch n {
	case 0:
		_ = NewSettings(&SSettings{})
	}
}

func TestInvalidMessage(t *testing.T) {
	t.Parallel()

	params := NewSettings(&SSettings{
		FMessageSizeBytes: (2 << 10),
	})

	if _, err := LoadMessage(params, struct{}{}); err == nil {
		t.Error("success load message with unknown type")
		return
	}

	if _, err := LoadMessage(params, []byte{123}); err == nil {
		t.Error("success load invalid message")
		return
	}

	if _, err := LoadMessage(params, []byte(CSeparator)); err == nil {
		t.Error("success unmarshal invalid message")
		return
	}

	if _, err := LoadMessage(params, tcInvalidPrefix+"!@#"); err == nil {
		t.Error("success decode body in invalid message")
		return
	}

	if _, err := LoadMessage(params, []byte(tcInvalidPrefix+"12")); err == nil {
		t.Error("success decode body with invalid message size")
		return
	}

	prng := random.NewStdPRNG()
	if _, err := LoadMessage(params, []byte(tcInvalidPrefix+"12"+prng.GetString(928))); err == nil {
		t.Error("success decode body with invalid message size")
		return
	}
}

func TestMessage(t *testing.T) {
	t.Parallel()

	params := NewSettings(&SSettings{
		FMessageSizeBytes: (2 << 10),
	})

	msg1, err := LoadMessage(params, testutils.TCBinaryMessage)
	if err != nil {
		t.Error(err)
		return
	}
	testMessage(t, msg1)

	msg2, err := LoadMessage(params, testutils.TCStringMessage)
	if err != nil {
		t.Error(err)
		return
	}
	testMessage(t, msg2)
}

func testMessage(t *testing.T, msg IMessage) {
	if !bytes.Equal(msg.GetSalt(), encoding.HexDecode(tcSalt)) {
		t.Error("incorrect salt value")
		return
	}

	if !bytes.Equal(msg.GetEncKey(), encoding.HexDecode(tcSession)) {
		t.Error("incorrect session value")
		return
	}

	if !bytes.Equal(msg.GetPubKey(), encoding.HexDecode(tcSender)) {
		t.Error("incorrect sender value")
		return
	}

	if !bytes.Equal(msg.GetSign(), encoding.HexDecode(tcSign)) {
		t.Error("incorrect sign value")
		return
	}

	if !bytes.Equal(msg.GetHash(), encoding.HexDecode(tcHash)) {
		t.Error("incorrect hash value")
		return
	}

	if msg.GetPayload() == nil {
		t.Error("failed get encrypted payload")
		return
	}

	msg1 := msg.(*SMessage)
	msg1.FPayload = []byte{123}
	if msg.GetPayload() != nil {
		t.Error("success got incorrect payload")
		return
	}
}
