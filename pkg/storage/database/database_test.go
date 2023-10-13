package database

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/storage"

	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcPathDBTemplate = "database_test_%d.db"
)

func TestTryDecrypt(t *testing.T) {
	cipher := symmetric.NewAESCipher([]byte("abcdefghijklmnopqrstuvwxyz123456"))
	if _, err := tryDecrypt(cipher, []byte{1}, []byte{2}); err == nil {
		t.Error("invalid size of encrypt data")
		return
	}

	authKey := []byte("auth-key")
	encData := cipher.EncryptBytes([]byte{})
	resData := bytes.Join(
		[][]byte{
			hashing.NewHMACSHA256Hasher(
				authKey,
				encData,
			).ToBytes(),
			encData,
		},
		[]byte{},
	)

	if _, err := tryDecrypt(cipher, authKey, resData); err != nil {
		t.Error("got error with encrypted empty slice: []")
		return
	}
}

func TestAllDBs(t *testing.T) {
	testCreate(t)
	testCipherKey(t)
	testBasic(t)
}

func testCreate(t *testing.T) {
	dbPath := fmt.Sprintf(tcPathDBTemplate, 3)
	defer os.RemoveAll(dbPath)

	store, err := NewKeyValueDB(storage.NewSettings(&storage.SSettings{
		FPath:     dbPath,
		FWorkSize: testutils.TCWorkSize,
		FPassword: "CIPHER",
	}))
	if err != nil {
		t.Error(err)
		return
	}
	defer store.Close()

	if store.GetSettings().GetPassword() != "CIPHER" {
		t.Error("[testCreate] incorrect default value = cipherKey")
		return
	}

	if store.GetSettings().GetPath() != dbPath {
		t.Error("[testCreate] incorrect default value = path")
		return
	}

	if err := store.Set([]byte("KEY"), []byte("VALUE")); err != nil {
		t.Error(err)
		return
	}

	if _, err := store.Get([]byte("KEY")); err != nil {
		t.Error(err)
		return
	}
}

func testCipherKey(t *testing.T) {
	dbPath := fmt.Sprintf(tcPathDBTemplate, 2)
	defer os.RemoveAll(dbPath)

	store, err := NewKeyValueDB(storage.NewSettings(&storage.SSettings{
		FPath:     dbPath,
		FWorkSize: testutils.TCWorkSize,
		FPassword: "CIPHER1",
	}))
	if err != nil {
		t.Error("[testCipherKey]", err)
		return
	}
	defer store.Close()

	store.Set([]byte("KEY"), []byte("VALUE"))

	store.Close() // open this database with another key
	store, err = NewKeyValueDB(storage.NewSettings(&storage.SSettings{
		FPath:     dbPath,
		FWorkSize: testutils.TCWorkSize,
		FPassword: "CIPHER2",
	}))
	if err != nil {
		t.Error("[testCipherKey]", err)
		return
	}
	defer store.Close()

	if _, err := store.Get([]byte("KEY")); err == nil {
		t.Error("[testCipherKey] success read value by another cipher key")
		return
	}
}

func testBasic(t *testing.T) {
	dbPath := fmt.Sprintf(tcPathDBTemplate, 1)
	defer os.RemoveAll(dbPath)

	store, err := NewKeyValueDB(storage.NewSettings(&storage.SSettings{
		FPath:     dbPath,
		FWorkSize: testutils.TCWorkSize,
		FPassword: "CIPHER",
	}))
	if err != nil {
		t.Error("[testBasic]", err)
		return
	}
	defer store.Close()

	secret1 := asymmetric.NewRSAPrivKey(512).ToBytes()
	store.Set([]byte("KEY"), secret1)

	secret2, err := store.Get([]byte("KEY"))
	if err != nil {
		t.Error("[testBasic]", err)
		return
	}

	if !bytes.Equal(secret1, secret2) {
		t.Error("[testBasic] saved and loaded values not equals")
		return
	}

	if err := store.Del([]byte("KEY")); err != nil {
		t.Error("[testBasic]", err)
		return
	}

	if _, err := store.Get([]byte("undefined key")); err == nil {
		t.Error("[testBasic] got value by undefined key")
		return
	}
}
