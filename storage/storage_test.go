package storage

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/settings/testutils"
)

func TestCryptoStorage(t *testing.T) {
	const (
		storageName = "storage.enc"
		storageKey  = "storage-key"

		key = "[application#1]password"
	)

	defer os.Remove(storageName)
	secret1 := crypto.NewPrivKey(512).Bytes()

	store := NewCryptoStorage(testutils.NewSettings(), storageName, []byte(storageKey))
	defer store.Close()

	store.Set([]byte(key), secret1)

	secret2, err := store.Get([]byte(key))
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(secret1, secret2) {
		t.Errorf("saved and loaded values not equals")
	}

	err = store.Del([]byte(key))
	if err != nil {
		t.Error(err)
	}
}

func TestLevelDBStorage(t *testing.T) {
	const (
		storageName = "storage.db"

		key = "[application#1]password"
	)

	defer os.RemoveAll(storageName)
	secret1 := crypto.NewPrivKey(512).Bytes()

	store := NewLevelDBStorage(storageName)
	defer store.Close()

	store.Set([]byte(key), secret1)

	secret2, err := store.Get([]byte(key))
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(secret1, secret2) {
		t.Errorf("saved and loaded values not equals")
	}

	err = store.Del([]byte(key))
	if err != nil {
		t.Error(err)
	}
}
