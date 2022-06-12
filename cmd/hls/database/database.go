package database

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/storage"
)

type sKeyValueDB struct {
	fMutex   sync.Mutex
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

func (db *sKeyValueDB) Push(key []byte) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	if len(key) != hashing.HashSize {
		return fmt.Errorf("hash size invalid")
	}

	if db.isExist(key) {
		return fmt.Errorf("hash already exists")
	}

	return db.fStorage.Set(getKeyHash(key), []byte{1})
}

func (db *sKeyValueDB) Exist(key []byte) bool {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	if len(key) != hashing.HashSize {
		return false
	}

	return db.isExist(key)
}

func (db *sKeyValueDB) Close() error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fStorage.Close()
}

func (db *sKeyValueDB) isExist(key []byte) bool {
	_, err := db.fStorage.Get(getKeyHash(key))
	return err == nil
}
