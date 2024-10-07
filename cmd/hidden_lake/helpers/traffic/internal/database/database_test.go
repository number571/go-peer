package database

import (
	"errors"
	"testing"

	"github.com/number571/go-peer/pkg/storage/database"
)

func TestVoidKVDatabase(t *testing.T) {
	db := NewVoidKVDatabase()
	if err := db.Set([]byte("aaa"), []byte("bbb")); err != nil {
		t.Error(err)
		return
	}
	_, err := db.Get([]byte("aaa"))
	if err == nil {
		t.Error("success get value from void database")
		return
	}
	if !errors.Is(err, database.ErrNotFound) {
		t.Error("got unsupported error")
		return
	}
	if err := db.Close(); err != nil {
		t.Error(err)
		return
	}
}
