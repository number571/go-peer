package database

import (
	"bytes"
	"fmt"
	"os"
	"testing"
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

func TestTryRecover(t *testing.T) {
	t.Parallel()

	dbPath := fmt.Sprintf(tcPathDBTemplate, 4)
	defer os.RemoveAll(dbPath)

	store, err := NewKVDatabase(dbPath)
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

func TestInvalidCreateDB(t *testing.T) {
	t.Parallel()

	path := "./not_exist/path/to/database/57199u140291724y121291d1/database.db"
	defer os.RemoveAll(path)

	_, err := NewKVDatabase(path)
	if err == nil {
		t.Error("success create database with incorrect path")
		return
	}
}

func TestCreateDB(t *testing.T) {
	t.Parallel()

	dbPath := fmt.Sprintf(tcPathDBTemplate, 3)
	defer os.RemoveAll(dbPath)

	store, err := NewKVDatabase(dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	defer store.Close()

	if err := store.Set([]byte("KEY"), []byte("VALUE")); err != nil {
		t.Error(err)
		return
	}

	if _, err := store.Get([]byte("KEY")); err != nil {
		t.Error(err)
		return
	}
}

func TestBasicDB(t *testing.T) {
	t.Parallel()

	dbPath := fmt.Sprintf(tcPathDBTemplate, 1)
	defer os.RemoveAll(dbPath)

	store, err := NewKVDatabase(dbPath)
	if err != nil {
		t.Error("[testBasic]", err)
		return
	}
	defer store.Close()

	data1 := []byte("hello, world!")
	if err := store.Set([]byte("KEY"), data1); err != nil {
		t.Error(err)
		return
	}

	data2, err := store.Get([]byte("KEY"))
	if err != nil {
		t.Error("[testBasic]", err)
		return
	}

	if !bytes.Equal(data1, data2) {
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
