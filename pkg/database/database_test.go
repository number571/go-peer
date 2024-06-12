package database

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/symmetric"

	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcPathDBTemplate = "database_test_%d.db"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SDatabaseError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestSettings(t *testing.T) {
	t.Parallel()

	for i := 0; i < 1; i++ {
		testSettings(t, i)
	}
}

func testSettings(t *testing.T, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	switch n { // nolint: gocritic
	case 0:
		_ = NewSettings(&SSettings{})
	}
}

func TestTryRecover(t *testing.T) {
	t.Parallel()

	dbPath := fmt.Sprintf(tcPathDBTemplate, 4)
	defer os.RemoveAll(dbPath)

	store, err := NewKVDatabase(NewSettings(&SSettings{
		FPath:     dbPath,
		FWorkSize: testutils.TCWorkSize,
		FPassword: "CIPHER",
	}))
	if err != nil {
		t.Error(err)
		return
	}

	if err := store.Set([]byte("KEY"), []byte("VALUE")); err != nil {
		t.Error(err)
		return
	}

	store.Close()
}

func TestTryDecrypt(t *testing.T) {
	t.Parallel()

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

	resData[0] ^= 1
	if _, err := tryDecrypt(cipher, authKey, resData); err == nil {
		t.Error("succes decrypt with corrupted data")
		return
	}
}

// // The test fails when the user is root.
// func TestInvalidCreateDB(t *testing.T) {
// 	t.Parallel()

// 	prng := random.NewCSPRNG()
// 	path := "/" + prng.GetString(32) + "/" + prng.GetString(32) + "/" + prng.GetString(32)
// 	defer os.RemoveAll(path)

// 	_, err := NewKVDatabase(database.NewSettings(&database.SSettings{
// 		FPath:     path,
// 		FWorkSize: testutils.TCWorkSize,
// 		FPassword: "CIPHER",
// 	}))
// 	if err == nil {
// 		t.Error("success create database with incorrect path")
// 		return
// 	}
// }

func TestInvalidInitDB(t *testing.T) {
	t.Parallel()

	testIvalidInitDB(t, 5, []byte(cHashKey))
	testIvalidInitDB(t, 6, []byte(cRandKey))
}

func testIvalidInitDB(t *testing.T, n int, key []byte) {
	dbPath := fmt.Sprintf(tcPathDBTemplate, n)
	defer os.RemoveAll(dbPath)

	store, err := NewKVDatabase(NewSettings(&SSettings{
		FPath:     dbPath,
		FWorkSize: testutils.TCWorkSize,
		FPassword: "CIPHER",
	}))
	if err != nil {
		t.Error(err)
		return
	}

	db := store.(*sKVDatabase)
	if err := delDB(db.fDB, key); err != nil {
		t.Error(err)
		return
	}

	store.Close()

	_, errx := NewKVDatabase(NewSettings(&SSettings{
		FPath:     dbPath,
		FWorkSize: testutils.TCWorkSize,
		FPassword: "CIPHER",
	}))
	if errx == nil {
		t.Error("success open database with incorrect param")
		return
	}
}

func TestCreateDB(t *testing.T) {
	t.Parallel()

	dbPath := fmt.Sprintf(tcPathDBTemplate, 3)
	defer os.RemoveAll(dbPath)

	store, err := NewKVDatabase(NewSettings(&SSettings{
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

func TestCipherKey(t *testing.T) {
	t.Parallel()

	dbPath := fmt.Sprintf(tcPathDBTemplate, 2)
	defer os.RemoveAll(dbPath)

	store, err := NewKVDatabase(NewSettings(&SSettings{
		FPath:     dbPath,
		FWorkSize: testutils.TCWorkSize,
		FPassword: "CIPHER1",
	}))
	if err != nil {
		t.Error("[testCipherKey]", err)
		return
	}
	defer store.Close()

	if err := store.Set([]byte("KEY"), []byte("VALUE")); err != nil {
		t.Error(err)
		return
	}

	store.Close() // open this database with another key
	_, err = NewKVDatabase(NewSettings(&SSettings{
		FPath:     dbPath,
		FWorkSize: testutils.TCWorkSize,
		FPassword: "CIPHER2",
	}))
	if err == nil {
		t.Error("[testCipherKey] success read database by another cipher key")
		return
	}
}

func TestBasicDB(t *testing.T) {
	t.Parallel()

	dbPath := fmt.Sprintf(tcPathDBTemplate, 1)
	defer os.RemoveAll(dbPath)

	store, err := NewKVDatabase(NewSettings(&SSettings{
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
	if err := store.Set([]byte("KEY"), secret1); err != nil {
		t.Error(err)
		return
	}

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
