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
		t.Fatal("incorrect err.Error()")
	}
}

func TestInvalidCreateDB(t *testing.T) {
	t.Parallel()

	path := "./not_exist/path/to/database/57199u140291724y121291d1/database.db"
	defer func() { _ = os.RemoveAll(path) }()

	_, err := NewKVDatabase(path)
	if err == nil {
		t.Fatal("success create database with incorrect path")
	}
}

func TestClosedDB(t *testing.T) {
	t.Parallel()

	dbPath := fmt.Sprintf(tcPathDBTemplate, 2)
	defer func() { _ = os.RemoveAll(dbPath) }()

	db, err := NewKVDatabase(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	if err := db.Set([]byte("KEY"), []byte("VALUE")); err == nil {
		t.Fatal("success set with closed db")
	}

	if err := db.Del([]byte("KEY")); err == nil {
		t.Fatal("success del with closed db")
	}
}

func TestCreateDB(t *testing.T) {
	t.Parallel()

	dbPath := fmt.Sprintf(tcPathDBTemplate, 3)
	defer func() { _ = os.RemoveAll(dbPath) }()

	store, err := NewKVDatabase(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = store.Close() }()

	if err := store.Set([]byte("KEY"), []byte("VALUE")); err != nil {
		t.Fatal(err)
	}

	if _, err := store.Get([]byte("KEY")); err != nil {
		t.Fatal(err)
	}

	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestBasicDB(t *testing.T) {
	t.Parallel()

	dbPath := fmt.Sprintf(tcPathDBTemplate, 1)
	defer func() { _ = os.RemoveAll(dbPath) }()

	store, err := NewKVDatabase(dbPath)
	if err != nil {
		t.Fatal("[testBasic]", err)
	}
	defer func() { _ = store.Close() }()

	if _, err := store.Get([]byte("KEY")); err == nil {
		t.Fatal("[testBasic] success get with bucket=nil")
	}

	if err := store.Del([]byte("KEY")); err != nil {
		t.Fatal("[testBasic]", err) // without error if bucket=nil
	}

	data1 := []byte("hello, world!")
	if err := store.Set([]byte("KEY"), data1); err != nil {
		t.Fatal(err)
	}

	data2, err := store.Get([]byte("KEY"))
	if err != nil {
		t.Fatal("[testBasic]", err)
	}

	if !bytes.Equal(data1, data2) {
		t.Fatal("[testBasic] saved and loaded values not equals")
	}

	if err := store.Del([]byte("KEY")); err != nil {
		t.Fatal("[testBasic]", err)
	}

	if _, err := store.Get([]byte("undefined key")); err == nil {
		t.Fatal("[testBasic] got value by undefined key")
	}
}
