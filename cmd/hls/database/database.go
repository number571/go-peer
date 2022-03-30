package database

import (
	"fmt"
	"os"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

type sKeyValueDB struct {
	fMutex sync.Mutex
	fPath  string
	fDB    *leveldb.DB
}

func NewKeyValueDB(path string) IKeyValueDB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		panic(err)
	}
	return &sKeyValueDB{
		fPath: path,
		fDB:   db,
	}
}

func (db *sKeyValueDB) Push(hash []byte) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	_, err := db.fDB.Get(GetKeyHash(hash), nil)
	if err == nil {
		return fmt.Errorf("hash already exists")
	}

	return db.fDB.Put(GetKeyHash(hash), []byte{1}, nil)
}

func (db *sKeyValueDB) Exist(hash []byte) bool {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	_, err := db.fDB.Get(GetKeyHash(hash), nil)
	return err == nil
}

func (db *sKeyValueDB) Close() error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fDB.Close()
}

func (db *sKeyValueDB) Clean() error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	db.fDB.Close()

	err := os.RemoveAll(db.fPath)
	if err != nil {
		return err
	}

	db.fDB = NewKeyValueDB(db.fPath).dbPointer()
	return nil
}

func (db *sKeyValueDB) dbPointer() *leveldb.DB {
	return db.fDB
}
