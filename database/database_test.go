package database

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/testutils"
)

const (
	tcPathDB    = "database_test.db"
	countOfIter = 10
)

func TestLevelDB(t *testing.T) {
	defer os.RemoveAll(tcPathDB)
	secret1 := asymmetric.NewRSAPrivKey(512).Bytes()

	store := NewLevelDB(&SSettings{
		FPath:      tcPathDB,
		FHashing:   true,
		FCipherKey: []byte(testutils.TcKey1),
	})
	defer store.Close()

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

func TestLevelDBIter(t *testing.T) {
	defer os.RemoveAll(tcPathDB)

	store := NewLevelDB(&SSettings{
		FPath:      tcPathDB,
		FHashing:   false,
		FCipherKey: []byte(testutils.TcKey1),
	})
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
	iter := store.Iter([]byte(testutils.TcKey2))
	defer iter.Close()

	for iter.Next() {
		val := string(iter.Value())
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
