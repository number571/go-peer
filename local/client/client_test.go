package client

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/local/payload"
	"github.com/number571/go-peer/local/routing"
	"github.com/number571/go-peer/settings"
)

const (
	tcHead = 0xDEADBEAF
	tcBody = "hello, world!"
)

func testNewClient() IClient {
	return NewClient(
		settings.NewSettings(),
		asymmetric.NewRSAPrivKey(1024),
	)
}

func TestEncrypt(t *testing.T) {
	client1 := testNewClient()
	client2 := testNewClient()

	pl := payload.NewPayload(uint64(tcHead), []byte(tcBody))
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

	if !bytes.Equal([]byte(tcBody), decPl.Body()) {
		t.Error("data not equal with decrypted data")
		return
	}
}
