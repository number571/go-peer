package storage

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	testStorageName = "storage.stg"
)

func TestCryptoStorage(t *testing.T) {
	os.Remove(testStorageName)
	defer os.Remove(testStorageName)

	sett := NewSettings(&SSettings{
		FPath:      testStorageName,
		FWorkSize:  testutils.TCWorkSize,
		FCipherKey: []byte("CIPHER"),
	})

	_, err := NewCryptoStorage(sett)
	if err != nil {
		t.Error(err)
		return
	}

	// try open already exists storage
	store, err := NewCryptoStorage(sett)
	if err != nil {
		t.Error(err)
		return
	}

	secret1 := asymmetric.NewRSAPrivKey(512).ToBytes()
	store.Set([]byte("KEY"), secret1)

	secret2, err := store.Get([]byte("KEY"))
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(secret1, secret2) {
		t.Error("saved and loaded values not equals")
		return
	}

	if err := store.Del([]byte("KEY")); err != nil {
		t.Error(err)
		return
	}

	if _, err := store.Get([]byte("KEY")); err == nil {
		t.Error("value in storage not deleted")
		return
	}
}
