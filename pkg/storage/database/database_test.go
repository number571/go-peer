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

const (
	tcPathDBTemplate = "database_test_%d.db"
	countOfIter      = 10
)

func TestFailCreateLevelDB(t *testing.T) {
	dbPath := fmt.Sprintf(
		"/fail_test_random/%s/111/222/333",
		random.NewStdPRNG().GetString(16),
	)
	defer os.RemoveAll(dbPath)

	store := NewLevelDB(NewSettings(&SSettings{
		FPath: dbPath,
	}))
	if store != nil {
		t.Errorf("this path '%s' realy exists?", dbPath)
		store.Close()
	}
}

func TestCreateLevelDB(t *testing.T) {
	defer os.RemoveAll(cPath)

	store := NewLevelDB(NewSettings(&SSettings{}))
	defer store.Close()

	if !bytes.Equal(store.GetSettings().GetCipherKey(), []byte(cCipherKey)) {
		t.Error("incorrect default value = cipherKey")
		return
	}

	if store.GetSettings().GetHashing() != false {
		t.Error("incorrect default value = hashing")
		return
	}

	if store.GetSettings().GetPath() != cPath {
		t.Error("incorrect default value = path")
		return
	}
}

func TestCipherKeyLevelDB(t *testing.T) {
	dbPath := fmt.Sprintf(tcPathDBTemplate, 3)
	defer os.RemoveAll(dbPath)

	store := NewLevelDB(NewSettings(&SSettings{
		FPath:      dbPath,
		FHashing:   true,
		FCipherKey: []byte(testutils.TcKey1),
	}))
	defer store.Close()

	store.Set([]byte(testutils.TcKey2), []byte(testutils.TcKey3))

	store.Close() // open this database with another key
	store = NewLevelDB(NewSettings(&SSettings{
		FPath:      dbPath,
		FHashing:   true,
		FCipherKey: []byte(testutils.TcKey2),
	}))
	defer store.Close()

	if _, err := store.Get([]byte(testutils.TcKey2)); err == nil {
		t.Error("success read value by another cipher key")
		return
	}
}

func TestLevelDB(t *testing.T) {
	dbPath := fmt.Sprintf(tcPathDBTemplate, 2)
	defer os.RemoveAll(dbPath)

	store := NewLevelDB(NewSettings(&SSettings{
		FPath:      dbPath,
		FHashing:   true,
		FCipherKey: []byte(testutils.TcKey1),
	}))
	defer store.Close()

	secret1 := asymmetric.NewRSAPrivKey(512).ToBytes()
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

	iter := store.GetIterator([]byte("_value"))
	if iter != nil {
		t.Error("iter is not null with hashing=true")
		return
	}

	if _, err := store.Get([]byte("undefined key")); err == nil {
		t.Error("got value by undefined key")
		return
	}
}

func TestLevelDBIter(t *testing.T) {
	dbPath := fmt.Sprintf(tcPathDBTemplate, 1)
	defer os.RemoveAll(dbPath)

	store := NewLevelDB(NewSettings(&SSettings{
		FPath:      dbPath,
		FHashing:   false,
		FCipherKey: []byte(testutils.TcKey1),
	}))
	defer store.Close()

	for i := 0; i < countOfIter; i++ {
		err := store.Set(
			[]byte(fmt.Sprintf("%s%d", testutils.TcKey2, i)),
			[]byte(fmt.Sprintf("%s%d", testutils.TcVal1, i)),
		)
		if err != nil {
			t.Error(err)
			return
		}
	}

	count := 0
	iter := store.GetIterator([]byte(testutils.TcKey2))
	defer iter.Close()

	for iter.Next() {
		if !bytes.Equal([]byte(fmt.Sprintf("%s%d", testutils.TcKey2, count)), iter.GetKey()) {
			t.Error("key not equal saved key")
			return
		}
		val := string(iter.GetValue())
		if val != fmt.Sprintf("%s%d", testutils.TcVal1, count) {
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
