package storage

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	tcStorageName    = "test_storage.stg"
	tcTmpStorageName = "test_storage-tmp.stg"

	tcSavedSecret = "3082013a020100024100cf3752dafd5eb10f00afa8a68c9b9853a91a56f2b4e6529f12ba953d835c3214233e4cc6e10d840e4388c7b84255d4ed73123028a121f288d85ad73b793a1887020301000102406ba34653e1075e1bf7f4473bf49022895aaf06f95e44c22845774c6cbe9e9697f5bf73bda74479ae793339ec31b3eaffd1fdf4ed7e8bc6794f6b85c25b983509022100db31184ed02bcba9bb0ca3f64df141b88d8866266ef10be46d8603808fa924e3022100f203672139df29f5699680e74b3593f0227d1c7410fd16b38879df5d43ae330d022026c369ef16358890fdb9608dc07ef80671513bef741340ed26c95a7933eecfcd022100e4d38ed97daca231a70a650b4cb37613a1a88614c0536cf987db23f53d1f22a902207ad51e4caf1bfa2b45899176780249dd9a5fa49803a98df1089b0564fb6dd140"
)

func TestCryptoStorage(t *testing.T) {
	testCryptoStorage(t, tcStorageName)
	testTempCryptoStorage(t, tcTmpStorageName)
}

func TestSettings(t *testing.T) {
	for i := 0; i < 3; i++ {
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
		_ = NewSettings(&SSettings{
			FWorkSize: testutils.TCWorkSize,
			FPassword: "CIPHER",
		})
	case 1:
		_ = NewSettings(&SSettings{
			FPath:     "path",
			FPassword: "CIPHER",
		})
	case 2:
		_ = NewSettings(&SSettings{
			FPath:     "path",
			FWorkSize: testutils.TCWorkSize,
		})
	}
}

func testCryptoStorage(t *testing.T, path string) {
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

func testTempCryptoStorage(t *testing.T, path string) {
	os.Remove(tcTmpStorageName)
	defer os.Remove(tcTmpStorageName)

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

	if _, err := store.Get([]byte("KEY")); err == nil {
		t.Error("value in storage not deleted")
		return
	}
}
