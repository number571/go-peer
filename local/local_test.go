package local

import (
	"bytes"
	"os"
	"testing"

	cr "github.com/number571/go-peer/crypto"
	gp "github.com/number571/go-peer/settings"
	tu "github.com/number571/go-peer/settings/testutils"
)

func newClient() Client {
	settings := tu.NewSettings()
	privKey := cr.NewPrivKey(settings.Get(gp.SizeAkey))
	return NewClient(privKey, settings)
}

func TestEncrypt(t *testing.T) {
	client1 := newClient()
	client2 := newClient()

	title := []byte("header")
	data := []byte("hello, world!")

	msg := NewMessage(title, data)
	encmsg, _ := client1.Encrypt(NewRoute(client2.PubKey(), nil, nil), msg)

	decmsg := client2.Decrypt(encmsg)
	title1, data1 := decmsg.Export()

	if !bytes.Equal(title, title1) {
		t.Errorf("title not equal with decrypted title")
	}

	if !bytes.Equal(data, data1) {
		t.Errorf("data not equal with decrypted data")
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
	secret1 := cr.NewPrivKey(512).Bytes()

	store := NewStorage(tu.NewSettings(), storageName, storagePasw)
	store.Write(subject, password, secret1)

	secret2, err := store.Read(subject, password)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(secret1, secret2) {
		t.Errorf("saved and loaded values not equals")
	}
}
