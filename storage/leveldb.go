package storage

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	_ IKeyValueStorage = &sLevelDBStorage{}
)

type sLevelDBStorage struct {
	fMutex sync.Mutex
	fDB    *leveldb.DB
}

func NewLevelDBStorage(path string) IKeyValueStorage {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil
	}
	return &sLevelDBStorage{
		fDB: db,
	}
}

func (db *sLevelDBStorage) Set(key []byte, value []byte) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fDB.Put(key, value, nil)
}

func (db *sLevelDBStorage) Get(key []byte) ([]byte, error) {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fDB.Get(key, nil)
}

func (db *sLevelDBStorage) Del(key []byte) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fDB.Delete(key, nil)
}

func (db *sLevelDBStorage) Close() error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fDB.Close()
}
