package storage

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/internal/testutils"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

const (
	storageName = "storage.stg"
)

func TestCryptoStorage(t *testing.T) {
	os.Remove(storageName)
	defer os.Remove(storageName)

	store, err := NewCryptoStorage(storageName, []byte(testutils.TcKey1), 10)
	if err != nil {
		t.Error(err)
		return
	}

	secret1 := asymmetric.NewRSAPrivKey(512).Bytes()
	store.Set([]byte(testutils.TcKey2), secret1)

	secret2, err := store.Get([]byte(testutils.TcKey2))
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(secret1, secret2) {
		t.Error("saved and loaded values not equals")
		return
	}

	if err := store.Del([]byte(testutils.TcKey2)); err != nil {
		t.Error(err)
		return
	}

	if _, err := store.Get([]byte(testutils.TcKey2)); err == nil {
		t.Errorf("value in storage not deleted")
		return
	}
}
