package storage

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/go-peer/test/utils"
)

const (
	tcStorageName         = "test_storage.stg"
	tcStorageNameTemplate = "test_storage-tmp-%d.stg"

	tcSavedSecret = "3082013a020100024100cf3752dafd5eb10f00afa8a68c9b9853a91a56f2b4e6529f12ba953d835c3214233e4cc6e10d840e4388c7b84255d4ed73123028a121f288d85ad73b793a1887020301000102406ba34653e1075e1bf7f4473bf49022895aaf06f95e44c22845774c6cbe9e9697f5bf73bda74479ae793339ec31b3eaffd1fdf4ed7e8bc6794f6b85c25b983509022100db31184ed02bcba9bb0ca3f64df141b88d8866266ef10be46d8603808fa924e3022100f203672139df29f5699680e74b3593f0227d1c7410fd16b38879df5d43ae330d022026c369ef16358890fdb9608dc07ef80671513bef741340ed26c95a7933eecfcd022100e4d38ed97daca231a70a650b4cb37613a1a88614c0536cf987db23f53d1f22a902207ad51e4caf1bfa2b45899176780249dd9a5fa49803a98df1089b0564fb6dd140"
)

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
	switch n {
	case 0:
		_ = NewSettings(&SSettings{})
	}
}

func TestCryptoStorage(t *testing.T) {
	t.Parallel()

	path := tcStorageName

	sett := NewSettings(&SSettings{
		FPath:     path,
		FWorkSize: testutils.TCWorkSize,
		FPassword: "CIPHER",
	})

	store, err := NewCryptoStorage(sett)
	if err != nil {
		t.Error(err)
		return
	}

	if store.GetSettings().GetPath() != path {
		t.Error("incorrect value from settings")
		return
	}

	// // used only once for set value into test storage
	// if err := store.Set([]byte("KEY"), encoding.HexDecode(tcSavedSecret)); err != nil {
	// 	t.Error(err)
	// 	return
	// }

	gotSecret, err := store.Get([]byte("KEY"))
	if err != nil {
		t.Error(err)
		return
	}

	if encoding.HexEncode(gotSecret) != tcSavedSecret {
		t.Error(err)
		return
	}
}

func TestTempCryptoStorage(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(tcStorageNameTemplate, 1)

	os.Remove(path)
	defer os.Remove(path)

	sett := NewSettings(&SSettings{
		FPath:     path,
		FWorkSize: testutils.TCWorkSize,
		FPassword: "CIPHER",
	})

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

	if err := store.Del([]byte("KEY")); err == nil {
		t.Error("success delete already deleted value")
		return
	}

	if _, err := store.Get([]byte("KEY")); err == nil {
		t.Error("value in storage not deleted")
		return
	}

	if err := store.Set([]byte("KEY1"), []byte("VALUE1")); err != nil {
		t.Error(err)
		return
	}

	os.Remove(path)

	if _, err := store.Get([]byte("KEY1")); err == nil {
		t.Error("success get value in not exist storage")
		return
	}

	if err := store.Del([]byte("KEY1")); err == nil {
		t.Error("success delete value in not exist storage")
		return
	}
}

func TestInvalidCreateCryptoStorage(t *testing.T) {
	t.Parallel()

	prng := random.NewStdPRNG()
	path := "/" + prng.GetString(32) + "/" + prng.GetString(32) + "/" + prng.GetString(32)

	os.Remove(path)
	defer os.Remove(path)

	sett := NewSettings(&SSettings{
		FPath:     path,
		FWorkSize: testutils.TCWorkSize,
		FPassword: "CIPHER",
	})

	_, err := NewCryptoStorage(sett)
	if err == nil {
		t.Error("success create storage with incorrect path")
		return
	}
}

func TestInvalidSizeCryptoStorage(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(tcStorageNameTemplate, 2)

	os.Remove(path)
	defer os.Remove(path)

	randBytes := random.NewStdPRNG().GetBytes(1)
	if err := os.WriteFile(path, randBytes, 0o644); err != nil {
		t.Error(err)
		return
	}

	sett := NewSettings(&SSettings{
		FPath:     path,
		FWorkSize: testutils.TCWorkSize,
		FPassword: "CIPHER",
	})

	_, err := NewCryptoStorage(sett)
	if err == nil {
		t.Error("success open storage with invalid structure")
		return
	}
}

func TestInvalidSetCryptoStorage(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(tcStorageNameTemplate, 3)

	os.Remove(path)
	defer os.Remove(path)

	randBytes := random.NewStdPRNG().GetBytes(128)
	if err := os.WriteFile(path, randBytes, 0o644); err != nil {
		t.Error(err)
		return
	}

	sett := NewSettings(&SSettings{
		FPath:     path,
		FWorkSize: testutils.TCWorkSize,
		FPassword: "CIPHER",
	})

	store, err := NewCryptoStorage(sett)
	if err != nil {
		t.Error(err)
		return
	}

	secret1 := asymmetric.NewRSAPrivKey(512).ToBytes()
	if err := store.Set([]byte("KEY"), secret1); err == nil {
		t.Error("success set with incorrect storage structure")
		return
	}

	os.Remove(path)

	if _, err := store.(*sCryptoStorage).decrypt(); err == nil {
		t.Error("success decrypt with deleted storage")
		return
	}
}

func TestInvalidDelCryptoStorage(t *testing.T) {
	t.Parallel()

	path := fmt.Sprintf(tcStorageNameTemplate, 4)

	os.Remove(path)
	defer os.Remove(path)

	sett := NewSettings(&SSettings{
		FPath:     path,
		FWorkSize: testutils.TCWorkSize,
		FPassword: "CIPHER",
	})

	store, err := NewCryptoStorage(sett)
	if err != nil {
		t.Error(err)
		return
	}

	if err := store.Set([]byte("KEY1"), []byte("VALUE1")); err != nil {
		t.Error(err)
		return
	}

	os.Remove(path)

	randBytes := random.NewStdPRNG().GetBytes(128)
	if err := os.WriteFile(path, randBytes, 0o644); err != nil {
		t.Error(err)
		return
	}

	if _, err := store.Get([]byte("KEY1")); err == nil {
		t.Error("success get value in corrupted storage")
		return
	}

	if err := store.Del([]byte("KEY1")); err == nil {
		t.Error("success delete value in corrupted storage")
		return
	}
}
