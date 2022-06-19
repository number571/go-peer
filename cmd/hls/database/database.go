package database

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/crypto/hashing"
	gp_database "github.com/number571/go-peer/database"
)

type sKeyValueDB struct {
	fMutex sync.Mutex
	fPath  string
	fDB    gp_database.IKeyValueDB
}

func NewKeyValueDB(path string) IKeyValueDB {
	db := gp_database.NewLevelDB(path)
	if db == nil {
		panic("storage is nil")
	}
	return &sKeyValueDB{
		fPath: path,
		fDB:   db,
	}
}

func (db *sKeyValueDB) Push(key []byte) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	if len(key) != hashing.GSHA256Size {
		return fmt.Errorf("hash size invalid")
	}

	if db.isExist(key) {
		return fmt.Errorf("hash already exists")
	}

	return db.fDB.Set(getKeyHash(key), []byte{1})
}

func (db *sKeyValueDB) Exist(key []byte) bool {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	if len(key) != hashing.GSHA256Size {
		return false
	}

	return db.isExist(key)
}

func (db *sKeyValueDB) Close() error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fDB.Close()
}

func (db *sKeyValueDB) isExist(key []byte) bool {
	_, err := db.fDB.Get(getKeyHash(key))
	return err == nil
}
