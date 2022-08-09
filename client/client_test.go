package client

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/payload"
	"github.com/number571/go-peer/routing"
	"github.com/number571/go-peer/testutils"
)

func testNewClient() IClient {
	return NewClient(
		NewSettings(10, (1<<10)),
		asymmetric.NewRSAPrivKey(1024),
	)
}

func TestEncrypt(t *testing.T) {
	client1 := testNewClient()
	client2 := testNewClient()

	pl := payload.NewPayload(uint64(testutils.TcHead), []byte(testutils.TcBody))
	msg := client1.Encrypt(routing.NewRoute(client2.PubKey()), pl)

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

	if !bytes.Equal([]byte(testutils.TcBody), decPl.Body()) {
		t.Error("data not equal with decrypted data")
		return
	}
}
