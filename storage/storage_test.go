package storage

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/utils/testutils"
)

const (
	storageName = "storage.stg"
)

func TestCryptoStorage(t *testing.T) {
	defer os.Remove(storageName)
	secret1 := asymmetric.NewRSAPrivKey(512).Bytes()

	store := NewCryptoStorage(storageName, []byte(testutils.TcKey1), 10)
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

	err = store.Del([]byte(testutils.TcKey2))
	if err != nil {
		t.Error(err)
		return
	}
}
