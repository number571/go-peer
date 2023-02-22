package storage

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	testStorageSize = 32
	testStorageName = "storage.stg"
)

func TestCryptoStorage(t *testing.T) {
	os.Remove(testStorageName)
	defer os.Remove(testStorageName)

	store, err := NewCryptoStorage(testStorageName, []byte(testutils.TcKey1), 10)
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
		t.Error("value in storage not deleted")
		return
	}
}

func TestMemoryStorage(t *testing.T) {
	store := NewMemoryStorage(testStorageSize)

	for i := 0; i < testStorageSize; i++ {
		store.Set([]byte(fmt.Sprintf("%d", i)), []byte{byte(i)})
	}

	for i := 0; i < testStorageSize; i++ {
		res, err := store.Get([]byte(fmt.Sprintf("%d", i)))
		if err != nil {
			t.Error(err)
			return
		}
		if !bytes.Equal([]byte{byte(i)}, res) {
			t.Error("values not equals (1)")
			return
		}
	}

	key := []byte("AAA-key")
	val := []byte{12, 23, 34, 45, 56}
	store.Set(key, val)

	if _, err := store.Get([]byte("0")); err == nil {
		t.Error("queue does not work")
		return
	}

	res, err := store.Get(key)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(val, res) {
		t.Error("values not equals (2)")
		return
	}

	if err := store.Del(key); err != nil {
		t.Error(err)
		return
	}

	if _, err := store.Get(key); err == nil {
		t.Error("delete does not work")
		return
	}

	if err := store.Del([]byte("undefined key")); err == nil {
		t.Error("success delete value by undefined key")
		return
	}
}
