package client

import (
	"bytes"
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/local/message"
	"github.com/number571/go-peer/local/routing"
	"github.com/number571/go-peer/settings/testutils"
)

func testNewClient() IClient {
	sett := testutils.NewSettings()
	privKey := asymmetric.NewRSAPrivKey(1024)
	return NewClient(privKey, sett)
}

func TestEncrypt(t *testing.T) {
	client1 := testNewClient()
	client2 := testNewClient()

	title := []byte("header")
	data := []byte("hello, world!")

	msg := message.NewMessage(title, data)
	encmsg, _ := client1.Encrypt(routing.NewRoute(client2.PubKey()), msg)

	decmsg, title1 := client2.Decrypt(encmsg)

	if !bytes.Equal(data, decmsg.Body().Data()) {
		t.Error("data not equal with decrypted data")
	}

	if !bytes.Equal(title, title1) {
		t.Error("title not equal with decrypted title")
	}
}
