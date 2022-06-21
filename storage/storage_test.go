package storage

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/settings"
)

func TestCryptoStorage(t *testing.T) {
	const (
		storageName = "storage.stg"
		storageKey  = "storage-key"
		objectKey   = "[application#1]password"
	)

	defer os.Remove(storageName)
	secret1 := asymmetric.NewRSAPrivKey(512).Bytes()

	store := NewCryptoStorage(settings.NewSettings(), storageName, []byte(storageKey))
	store.Set([]byte(objectKey), secret1)

	secret2, err := store.Get([]byte(objectKey))
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(secret1, secret2) {
		t.Error("saved and loaded values not equals")
		return
	}

	err = store.Del([]byte(objectKey))
	if err != nil {
		t.Error(err)
		return
	}
}
