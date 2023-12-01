package client

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcMessageSize = (2 << 10)
	tcKeySizeBits = 1024
)

var (
	tgMsgLimit = testNewClient().GetMessageLimit()
	tgPrivKey  = asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024)
	tgMessages = []string{
		testutils.TcBody,
		"",
		"A",
		"AA",
		"AAA",
		"AAAA",
		"AAAAA",
		"AAAAAA",
		"AAAAAAA",
		"AAAAAAAA",
		"AAAAAAAAA",
		"AAAAAAAAAA",
		"AAAAAAAAAAA",
		"AAAAAAAAAAAA",
		"AAAAAAAAAAAAA",
		"AAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAAAAA",
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		random.NewStdPRNG().GetString(tgMsgLimit), // maximum size of message
	}
)

func TestClientWithMessageSize(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	_ = NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: 1024,
			FKeySizeBits:      tcKeySizeBits,
		}),
		tgPrivKey,
	)
}

func TestEncrypt(t *testing.T) {
	t.Parallel()

	client1 := testNewClient()
	client2 := testNewClient()

	_ = client1.GetSettings()
	_ = client1.GetPrivKey()

	pl := payload.NewPayload(uint64(testutils.TcHead), []byte(testutils.TcBody))
	msg, err := client1.EncryptPayload(client2.GetPubKey(), pl)
	if err != nil {
		t.Error(err)
		return
	}

	// os.WriteFile("test_binary.msg", msg.ToBytes(), 0644)
	// os.WriteFile("test_string.msg", []byte(msg.ToString()), 0644)

	_, decPl, err := client2.DecryptMessage(msg)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal([]byte(testutils.TcBody), decPl.GetBody()) {
		t.Error("data not equal with decrypted data")
		return
	}
}

func TestDecrypt(t *testing.T) {
	t.Parallel()

	client1 := testNewClient()

	pl := payload.NewPayload(uint64(testutils.TcHead), []byte(testutils.TcBody))
	msg, err := client1.EncryptPayload(client1.GetPubKey(), pl)
	if err != nil {
		t.Error(err)
		return
	}

	if _, _, err := client1.DecryptMessage(nil); err == nil {
		t.Error("success decrypt nil message")
		return
	}

	sMsg := msg.(*message.SMessage)

	sMsg1 := *sMsg
	sMsg1.FHash = "0"
	if _, _, err := client1.DecryptMessage(&sMsg1); err == nil {
		t.Error("success decrypt message with incorrect hash")
		return
	}

	sMsg3 := *sMsg
	sMsg3.FEncKey = "0"
	if _, _, err := client1.DecryptMessage(&sMsg3); err == nil {
		t.Error("success decrypt message with incorrect session key")
		return
	}

	sMsg4 := *sMsg
	sMsg4.FPubKey = "0"
	if _, _, err := client1.DecryptMessage(&sMsg4); err == nil {
		t.Error("success decrypt message with incorrect sender key (iv block)")
		return
	}

	sMsg5 := *sMsg
	sMsg5.FPubKey = "11111111111111111111111111111111"
	if _, _, err := client1.DecryptMessage(&sMsg5); err == nil {
		t.Error("success decrypt message with incorrect sender key (public key is nil)")
		return
	}

	sMsg6 := *sMsg
	sMsg6.FPayload = []byte{111}
	if _, _, err := client1.DecryptMessage(&sMsg6); err == nil {
		t.Error("success decrypt message with incorrect payload (iv block)")
		return
	}

	sMsg7 := *sMsg
	sMsg7.FPayload = []byte("11111111111111111111111111111111")
	if _, _, err := client1.DecryptMessage(&sMsg7); err == nil {
		t.Error("success decrypt message with incorrect payload (payload is nil)")
		return
	}

	sMsg8 := *sMsg
	sMsg8.FSalt = "0"
	if _, _, err := client1.DecryptMessage(&sMsg8); err == nil {
		t.Error("success decrypt message with incorrect salt")
		return
	}

	sMsg9 := *sMsg
	sMsg9.FHash = "111da32433c2d7f99b38042d7b73db291bd803c55f3c83745ae3ebae6ba111"
	if _, _, err := client1.DecryptMessage(&sMsg9); err == nil {
		t.Error("success decrypt message with incorrect hash check")
		return
	}

	sMsg10 := *sMsg
	sMsg10.FSign = "0"
	if _, _, err := client1.DecryptMessage(&sMsg10); err == nil {
		t.Error("success decrypt message with incorrect sign")
		return
	}

	sMsg11 := *sMsg
	sMsg11.FSign = "111ce5d111c74a0b8638f24f8ff200f64ca0e88cda1fd483783930b08e465fa9fc9565a0a3afbdfdf3f463bc77e526f2c41c6ddd2dae5d6f90e741442e2939731cbdad4071c29eff83dff932589b2cbfd8fa8a5fac19de4c40c3adde4cde1235c0bbf053b0e04e826993f8060a50c671c6bf56ce24fe4e921b60f6ca2239932ebd1b8c8556d5a2ac13e5ef1d8ea9c111"
	if _, _, err := client1.DecryptMessage(&sMsg11); err == nil {
		t.Error("success decrypt message with incorrect sign check")
		return
	}
}

func TestMessageSize(t *testing.T) {
	t.Parallel()

	client1 := testNewClient()
	sizes := make([]int, 0, len(tgMessages))

	for _, smsg := range tgMessages {
		pl := payload.NewPayload(uint64(testutils.TcHead), []byte(smsg))
		msg, err := client1.EncryptPayload(client1.GetPubKey(), pl)
		if err != nil {
			t.Error(err)
			return
		}
		sizes = append(sizes, len(msg.ToBytes()))
	}

	for i := 0; i < len(sizes)-1; i++ {
		if sizes[i] != sizes[i+1] {
			t.Errorf(
				"len bytes of different messages = id(%d, %d) not equals = size(%d, %d)",
				i, i+1,
				sizes[i], sizes[i+1],
			)
			return
		}
	}
}

func TestGetMessageLimit(t *testing.T) {
	t.Parallel()

	client1 := testNewClient()

	msg1 := random.NewStdPRNG().GetBytes(tgMsgLimit)
	pld1 := payload.NewPayload(uint64(testutils.TcHead), []byte(msg1))
	if _, err := client1.EncryptPayload(client1.GetPubKey(), pld1); err != nil {
		t.Error("message1 > message limit:", err)
		return
	}

	msg2 := random.NewStdPRNG().GetBytes(tgMsgLimit + 1)
	pld2 := payload.NewPayload(uint64(testutils.TcHead), []byte(msg2))
	if _, err := client1.EncryptPayload(client1.GetPubKey(), pld2); err == nil {
		t.Error("message2 > message limit but not alert:", err)
		return
	}
}

func testNewClient() IClient {
	return NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: tcMessageSize,
			FKeySizeBits:      tcKeySizeBits,
		}),
		tgPrivKey,
	)
}
