package local

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/settings/testutils"
)

func newClient() IClient {
	sett := testutils.NewSettings()
	privKey := crypto.NewPrivKey(1024)
	return NewClient(privKey, sett)
}

func TestEncrypt(t *testing.T) {
	client1 := newClient()
	client2 := newClient()

	title := []byte("header")
	data := []byte("hello, world!")

	msg := NewMessage(title, data)
	encmsg, _ := client1.Encrypt(NewRoute(client2.PubKey(), nil, nil), msg)

	decmsg, title1 := client2.Decrypt(encmsg)

	if !bytes.Equal(data, decmsg.Body().Data()) {
		t.Errorf("data not equal with decrypted data")
	}

	if !bytes.Equal(title, title1) {
		t.Errorf("title not equal with decrypted title")
	}
}

func TestStorage(t *testing.T) {
	const (
		storageName = "storage.enc"
		storagePasw = "storage-password"

		subject  = "application_#1"
		password = "privkey-password"
	)

	defer os.Remove(storageName)
	secret1 := crypto.NewPrivKey(512).Bytes()

	store := NewStorage(testutils.NewSettings(), storageName, storagePasw)
	store.Write(subject, password, secret1)

	secret2, err := store.Read(subject, password)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(secret1, secret2) {
		t.Errorf("saved and loaded values not equals")
	}
}
