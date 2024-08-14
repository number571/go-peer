// nolint: goerr113
package client

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcMessageSize = (2 << 10)
)

var (
	tgMsgLimit = testNewClient().GetMessageLimit()
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
		random.NewCSPRNG().GetString(tgMsgLimit - encoding.CSizeUint64), // maximum size of message - payload64.head
	}
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SClientError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestClientPanicWithMessageSize(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	_ = NewClient(
		NewSettings(&SSettings{}),
	)
}

func TestEncrypt(t *testing.T) {
	t.Parallel()

	client := testNewClient()
	key := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456")

	pl := payload.NewPayload64(uint64(testutils.TcHead), []byte(testutils.TcBody))
	msg, err := client.EncryptMessage(key, pl.ToBytes())
	if err != nil {
		t.Error(err)
		return
	}

	// lMsg, _ := message.LoadMessage(client1.GetSettings(), msg)
	// os.WriteFile("test_binary.msg", lMsg.ToBytes(), 0644)
	// os.WriteFile("test_string.msg", []byte(lMsg.ToString()), 0644)

	decMsg, err := client.DecryptMessage(key, msg)
	if err != nil {
		t.Error(err)
		return
	}

	decPl := payload.LoadPayload64(decMsg)
	if !bytes.Equal([]byte(testutils.TcBody), decPl.GetBody()) {
		t.Error("data not equal with decrypted data")
		return
	}
}

func TestMessageSize(t *testing.T) {
	t.Parallel()

	client1 := testNewClient()
	key := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456")

	for _, smsg := range tgMessages {
		pl := payload.NewPayload64(uint64(testutils.TcHead), []byte(smsg))
		msg, err := client1.EncryptMessage(key, pl.ToBytes())
		if err != nil {
			t.Error(err)
			return
		}
		if uint64(len(msg)) != client1.GetSettings().GetMessageSizeBytes() {
			t.Error("got invalid message size bytes")
			return
		}
	}
}

func TestGetMessageLimit(t *testing.T) {
	t.Parallel()

	client := testNewClient()
	key := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456")

	msg1 := random.NewCSPRNG().GetBytes(tgMsgLimit - encoding.CSizeUint64)
	pld1 := payload.NewPayload64(uint64(testutils.TcHead), msg1)
	if _, err := client.EncryptMessage(key, pld1.ToBytes()); err != nil {
		t.Error("message1 > message limit:", err)
		return
	}

	msg2 := random.NewCSPRNG().GetBytes(tgMsgLimit + 1)
	pld2 := payload.NewPayload64(uint64(testutils.TcHead), msg2)
	if _, err := client.EncryptMessage(key, pld2.ToBytes()); err == nil {
		t.Error("message2 > message limit but not alert:", err)
		return
	}
}

func testNewClient() IClient {
	return NewClient(
		NewSettings(&SSettings{
			FMessageSizeBytes: tcMessageSize,
		}),
	)
}
