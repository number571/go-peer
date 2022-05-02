package database

import (
	"fmt"
	"os"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/storage"
)

type sKeyValueDB struct {
	fPath    string
	fStorage storage.IKeyValueStorage
}

func NewKeyValueDB(path string) IKeyValueDB {
	stg := storage.NewLevelDBStorage(path)
	if stg == nil {
		panic("storage is nil")
	}
	return &sKeyValueDB{
		fPath:    path,
		fStorage: stg,
	}
}

func (db *sKeyValueDB) Push(hash []byte) error {
	if len(hash) != crypto.HashSize {
		return fmt.Errorf("hash size invalid")
	}
	if db.Exist(hash) {
		return fmt.Errorf("hash already exists")
	}
	return db.fStorage.Set(getKeyHash(hash), []byte{1})
}

func (db *sKeyValueDB) Exist(hash []byte) bool {
	if len(hash) != crypto.HashSize {
		return false
	}
	_, err := db.fStorage.Get(getKeyHash(hash))
	return err == nil
}

func (db *sKeyValueDB) Close() error {
	return db.fStorage.Close()
}

func (db *sKeyValueDB) Clean() error {
	db.fStorage.Close()

	err := os.RemoveAll(db.fPath)
	if err != nil {
		return err
	}

	db.fStorage = storage.NewLevelDBStorage(db.fPath)
	if db.fStorage == nil {
		panic("storage is nil")
	}

	return nil
}
