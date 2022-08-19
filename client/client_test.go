package client

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/payload"
	"github.com/number571/go-peer/testutils"
)

var (
	// AES block size = 16 bytes
	tgMessages = []string{
		testutils.TcBody,
		testutils.TcLargeBody,
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
	}
)

func testNewClient() IClient {
	return NewClient(
		NewSettings(10, (1<<20)),
		asymmetric.NewRSAPrivKey(1024),
	)
}

func TestEncrypt(t *testing.T) {
	client1 := testNewClient()
	client2 := testNewClient()

	pl := payload.NewPayload(uint64(testutils.TcHead), []byte(testutils.TcLargeBody))
	msg := client1.Encrypt(client2.PubKey(), pl)

	msgBytes := msg.Bytes()

	_, decPl := client2.Decrypt(msg)
	if decPl == nil {
		t.Error("decrypt payload is nil")
		return
	}

	if !bytes.Equal(msgBytes, msg.Bytes()) {
		t.Error("encrypted bytes not equal after action")
		return
	}

	if !bytes.Equal([]byte(testutils.TcLargeBody), decPl.Body()) {
		t.Error("data not equal with decrypted data")
		return
	}
}

func TestMessageSize(t *testing.T) {
	client1 := testNewClient()
	sizes := make([]int, 0, len(tgMessages))

	for _, smsg := range tgMessages {
		pl := payload.NewPayload(uint64(testutils.TcHead), []byte(smsg))
		msg := client1.Encrypt(client1.PubKey(), pl)
		sizes = append(sizes, len(msg.Bytes()))
	}

	for i := 0; i < len(sizes)-1; i++ {
		if sizes[i] != sizes[i+1] {
			t.Errorf("len bytes of different messages (%d, %d) not equals", i, i+1)
			return
		}
	}
}
