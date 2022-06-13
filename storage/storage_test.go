package storage

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/settings/testutils"
)

func TestCryptoStorage(t *testing.T) {
	const (
		storageName = "storage.stg"
		storageKey  = "storage-key"
		objectKey   = "[application#1]password"
	)

	defer os.Remove(storageName)
	secret1 := asymmetric.NewRSAPrivKey(512).Bytes()

	store := NewCryptoStorage(testutils.NewSettings(), storageName, []byte(storageKey))
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

func TestLevelDBStorage(t *testing.T) {
	const (
		storageName   = "storage.db"
		objectKey     = "[application#1]password"
		objectKeyIter = "object-key-"
		objectValIter = "object-value-"
		countOfIter   = 10
	)

	defer os.RemoveAll(storageName)
	secret1 := asymmetric.NewRSAPrivKey(512).Bytes()

	store := NewLevelDBStorage(storageName)
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
