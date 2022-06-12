package storage

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	_ IKeyValueStorage = &sLevelDBStorage{}
	_ iIterator        = &sLevelDBIterator{}
)

type sLevelDBStorage struct {
	fMutex sync.Mutex
	fDB    *leveldb.DB
}

type sLevelDBIterator struct {
	fMutex sync.Mutex
	ptr    iterator.Iterator
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

func (db *sLevelDBStorage) Iter(prefix []byte) iIterator {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return &sLevelDBIterator{
		ptr: db.fDB.NewIterator(util.BytesPrefix(prefix), nil),
	}
}

func (iter *sLevelDBIterator) Next() bool {
	iter.fMutex.Lock()
	defer iter.fMutex.Unlock()

	return iter.ptr.Next()
}

func (iter *sLevelDBIterator) Key() []byte {
	iter.fMutex.Lock()
	defer iter.fMutex.Unlock()

	return iter.ptr.Key()
}

func (iter *sLevelDBIterator) Value() []byte {
	iter.fMutex.Lock()
	defer iter.fMutex.Unlock()

	return iter.ptr.Value()
}

func (iter *sLevelDBIterator) Close() {
	iter.fMutex.Lock()
	defer iter.fMutex.Unlock()

	iter.ptr.Release()
}
