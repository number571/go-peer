package message

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcSalt    = "d84fc309a060cf101d040bc76a40e532171cc03fa07c9cbba82f0fcd96d14de37c8ccac1435a95c1826c1df58d7ce1ba"
	tcSession = "769692f8587e23cd184e8366487a400321cb5c5ca7561a3cb0f93efdcf8d0564bc6a47679ff2c0bc4ac50dacf92abad7dde6db7980153b2be5dbd5fce90f4e0836e3ffd5540842a4b25538b6f404fa51a38010df87f212c64549e819d5b2610acf0ceaf163018a74468c2bf0f190cbabdaef3dd4edc9175adebc3b69121b0971"
	tcSender  = "5ce9454a7fb51b47087821b4833bce2984cc730b832461919a0d063a7b66ea587f8c7a3e058a53119a84d2a9653a40285161de843a031ad064330ab942baf7f1a79400dd17f9ad5a77efa969c31704a6990a550cfe2eb0435aa77e3abaf5b2a0006944e2d129e2504e6edf2a17a8ba4376a1d19ed92320976e0b131d992429a32bea8f7f00c7b83acf5a0f3315fef483886efd10d77d1ad46e94d354"
	tcHash    = "8261038e2fcf3ea2338803d870f1ef7c802b0f52e7309bcf2dcf0d44721ee66f31618aad7cfb3d1669fe9f91b1959d95"
	tcSign    = "9d1a39e112472d42f0ecaada7e503a43dfcbb4b696388b54ed95d96d76893cd5dbc0236e3e8025ff39aa83a15f6ac91455a81bdedb319037a755afd1c6eb321d87871ee1492dd1363b3643fca2ea26739bc16a94d1fd64bfadff32ffb764117800e17c8d4d651df9339785e8981341408f39cf2af88aee4225a1f1249bf3faaf62c08f106f91298f4ba865d71041a208"

	// -1 byte from hash
	tcInvalidPrefix = `{"pubk":"5ce9454a7fb51b47087821b4833bce2984cc730b832461919a0d063a7b66ea587f8c7a3e058a53119a84d2a9653a40285161de843a031ad064330ab942baf7f1a79400dd17f9ad5a77efa969c31704a6990a550cfe2eb0435aa77e3abaf5b2a0006944e2d129e2504e6edf2a17a8ba4376a1d19ed92320976e0b131d992429a32bea8f7f00c7b83acf5a0f3315fef483886efd10d77d1ad46e94d354","enck":"769692f8587e23cd184e8366487a400321cb5c5ca7561a3cb0f93efdcf8d0564bc6a47679ff2c0bc4ac50dacf92abad7dde6db7980153b2be5dbd5fce90f4e0836e3ffd5540842a4b25538b6f404fa51a38010df87f212c64549e819d5b2610acf0ceaf163018a74468c2bf0f190cbabdaef3dd4edc9175adebc3b69121b0971","salt":"d84fc309a060cf101d040bc76a40e532171cc03fa07c9cbba82f0fcd96d14de37c8ccac1435a95c1826c1df58d7ce1ba","hash":"8261038e2fcfa2338803d870f1ef7c802b0f52e7309bcf2dcf0d44721ee66f31618aad7cfb3d1669fe9f91b1959d95","sign":"9d1a39e112472d42f0ecaada7e503a43dfcbb4b696388b54ed95d96d76893cd5dbc0236e3e8025ff39aa83a15f6ac91455a81bdedb319037a755afd1c6eb321d87871ee1492dd1363b3643fca2ea26739bc16a94d1fd64bfadff32ffb764117800e17c8d4d651df9339785e8981341408f39cf2af88aee4225a1f1249bf3faaf62c08f106f91298f4ba865d71041a208"}@`
)

func TestSettings(t *testing.T) {
	t.Parallel()

	for i := 0; i < 2; i++ {
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
		_ = NewSettings(&SSettings{
			FMessageSizeBytes: 1024,
		})
	case 1:
		_ = NewSettings(&SSettings{
			FKeySizeBits: testutils.TcKeySize,
		})
	}
}

func TestInvalidMessage(t *testing.T) {
	t.Parallel()

	params := NewSettings(&SSettings{
		FMessageSizeBytes: (2 << 10),
		FKeySizeBits:      testutils.TcKeySize,
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
		FKeySizeBits:      testutils.TcKeySize,
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
	if !bytes.Equal(msg.GetPayload(), []byte{123}) {
		t.Error("success got incorrect payload")
		return
	}
}
