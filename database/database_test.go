package database

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
)

func TestLevelDB(t *testing.T) {
	const (
		storageName = "storage.db"
		storageKey  = "storage-key"
		objectKey   = "[application#1]password"
	)

	defer os.RemoveAll(storageName)
	secret1 := asymmetric.NewRSAPrivKey(512).Bytes()

	store := NewLevelDB(storageName).
		WithHashing(true).
		WithEncryption([]byte(storageKey))
	defer store.Close()

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

func TestLevelDBIter(t *testing.T) {
	const (
		storageName   = "storage.db"
		storageKey    = "storage-key"
		objectKeyIter = "object-key-"
		objectValIter = "object-value-"
		countOfIter   = 10
	)

	defer os.RemoveAll(storageName)

	store := NewLevelDB(storageName).
		WithEncryption([]byte(storageKey))
	defer store.Close()

	for i := 0; i < countOfIter; i++ {
		err := store.Set(
			[]byte(fmt.Sprintf("%s%d", objectKeyIter, i)),
			[]byte(fmt.Sprintf("%s%d", objectValIter, i)),
		)
		if err != nil {
			t.Error(err)
			return
		}
	}

	count := 0
	iter := store.Iter([]byte(objectKeyIter))
	defer iter.Close()

	for iter.Next() {
		val := string(iter.Value())
		if val != fmt.Sprintf("%s%d", objectValIter, count) {
			t.Error("value not equal saved value")
			return
		}
		count++
	}

	if count != countOfIter {
		t.Error("count not equal count of iter")
		return
	}
}
