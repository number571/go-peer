package database

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

type (
	tiDBConsctruct func(ISettings) (IKVDatabase, error)
)

const (
	tcPathDBTemplate = "database_test_%d.db"
)

func TestAllDBs(t *testing.T) {
	testCreate(t, NewKeyValueDB)
	testCipherKey(t, NewKeyValueDB)
	testBasic(t, NewKeyValueDB)
}

func testCreate(t *testing.T, dbConstruct tiDBConsctruct) {
	dbPath := fmt.Sprintf(tcPathDBTemplate, 3)
	defer os.RemoveAll(dbPath)

	store, err := dbConstruct(NewSettings(&SSettings{
		FPath:      dbPath,
		FHashing:   false,
		FCipherKey: []byte("CIPHER"),
	}))
	if err != nil {
		t.Error(err)
		return
	}
	defer store.Close()

	if !bytes.Equal(store.GetSettings().GetCipherKey(), []byte("CIPHER")) {
		t.Error("[testCreate] incorrect default value = cipherKey")
		return
	}

	if store.GetSettings().GetHashing() != false {
		t.Error("[testCreate] incorrect default value = hashing")
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

func testCipherKey(t *testing.T, dbConstruct tiDBConsctruct) {
	dbPath := fmt.Sprintf(tcPathDBTemplate, 2)
	defer os.RemoveAll(dbPath)

	store, err := dbConstruct(NewSettings(&SSettings{
		FPath:      dbPath,
		FHashing:   true,
		FCipherKey: []byte("CIPHER1"),
	}))
	if err != nil {
		t.Error("[testCipherKey]", err)
		return
	}
	defer store.Close()

	store.Set([]byte("KEY"), []byte("VALUE"))

	store.Close() // open this database with another key
	store, err = dbConstruct(NewSettings(&SSettings{
		FPath:      dbPath,
		FHashing:   true,
		FCipherKey: []byte("CIPHER2"),
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

func testBasic(t *testing.T, dbConstruct tiDBConsctruct) {
	dbPath := fmt.Sprintf(tcPathDBTemplate, 1)
	defer os.RemoveAll(dbPath)

	store, err := dbConstruct(NewSettings(&SSettings{
		FPath:      dbPath,
		FHashing:   true,
		FCipherKey: []byte("CIPHER"),
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

	err = store.Del([]byte("KEY"))
	if err != nil {
		t.Error("[testBasic]", err)
		return
	}

	if _, err := store.Get([]byte("undefined key")); err == nil {
		t.Error("[testBasic] got value by undefined key")
		return
	}
}
