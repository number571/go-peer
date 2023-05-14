package database

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	testutils "github.com/number571/go-peer/test/_data"
)

type (
	tiDBConsctruct func(ISettings) (IKeyValueDB, error)
)

const (
	tcPathDBTemplate = "database_test_%d.db"
)

func TestAllDBs(t *testing.T) {
	testFailCreate(t, NewSQLiteDB)
	testCreate(t, NewSQLiteDB)
	testCipherKey(t, NewSQLiteDB)
	testBasic(t, NewSQLiteDB)
}

func testFailCreate(t *testing.T, dbConstruct tiDBConsctruct) {
	dbPath := fmt.Sprintf(
		"/fail_test_random/%s/111/222/333",
		random.NewStdPRNG().GetString(16),
	)
	defer os.RemoveAll(dbPath)

	store, err := dbConstruct(NewSettings(&SSettings{
		FPath: dbPath,
	}))
	if err == nil {
		t.Errorf("[testFailCreate] incorrect: error is nil")
		return
	}
	if store != nil {
		t.Errorf("[testFailCreate] this path '%s' realy exists?", dbPath)
		store.Close()
	}
}

func testCreate(t *testing.T, dbConstruct tiDBConsctruct) {
	defer os.RemoveAll(cPath)

	store, err := dbConstruct(NewSettings(&SSettings{}))
	if err != nil {
		t.Error(err)
		return
	}
	defer store.Close()

	if !bytes.Equal(store.GetSettings().GetCipherKey(), []byte(cCipherKey)) {
		t.Error("[testCreate] incorrect default value = cipherKey")
		return
	}

	if store.GetSettings().GetHashing() != false {
		t.Error("[testCreate] incorrect default value = hashing")
		return
	}

	if store.GetSettings().GetPath() != cPath {
		t.Error("[testCreate] incorrect default value = path")
		return
	}

	if err := store.Set([]byte(testutils.TcKey2), []byte(testutils.TcKey3)); err != nil {
		t.Error(err)
		return
	}

	if _, err := store.Get([]byte(testutils.TcKey2)); err != nil {
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
		FCipherKey: []byte(testutils.TcKey1),
	}))
	if err != nil {
		t.Error("[testCipherKey]", err)
		return
	}
	defer store.Close()

	store.Set([]byte(testutils.TcKey2), []byte(testutils.TcKey3))

	store.Close() // open this database with another key
	store, err = dbConstruct(NewSettings(&SSettings{
		FPath:      dbPath,
		FHashing:   true,
		FCipherKey: []byte(testutils.TcKey2),
	}))
	if err != nil {
		t.Error("[testCipherKey]", err)
		return
	}
	defer store.Close()

	if _, err := store.Get([]byte(testutils.TcKey2)); err == nil {
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
		FCipherKey: []byte(testutils.TcKey1),
	}))
	if err != nil {
		t.Error("[testBasic]", err)
		return
	}
	defer store.Close()

	secret1 := asymmetric.NewRSAPrivKey(512).ToBytes()
	store.Set([]byte(testutils.TcKey2), secret1)

	secret2, err := store.Get([]byte(testutils.TcKey2))
	if err != nil {
		t.Error("[testBasic]", err)
		return
	}

	if !bytes.Equal(secret1, secret2) {
		t.Error("[testBasic] saved and loaded values not equals")
		return
	}

	err = store.Del([]byte(testutils.TcKey2))
	if err != nil {
		t.Error("[testBasic]", err)
		return
	}

	if _, err := store.Get([]byte("undefined key")); err == nil {
		t.Error("[testBasic] got value by undefined key")
		return
	}
}
