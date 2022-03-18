package database

import (
	"fmt"
	"os"
	"sync"

	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/local"
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

func (db *sKeyValueDB) Size(key []byte) uint64 {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	data, err := db.fDB.Get(GetKeySize(key), nil)
	if err != nil {
		return 0
	}

	return encoding.BytesToUint64(data)
}

func (db *sKeyValueDB) Push(key []byte, msg local.IMessage) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	// store hash
	hash := msg.Body().Hash()
	_, err := db.fDB.Get(GetKeyHash(hash), nil)
	if err == nil {
		return fmt.Errorf("hash already exists")
	}

	err = db.fDB.Put(GetKeyHash(hash), []byte{1}, nil)
	if err != nil {
		return err
	}

	// update size
	size := uint64(0)
	bnum, err := db.fDB.Get(GetKeySize(key), nil)
	if err == nil {
		size = encoding.BytesToUint64(bnum)
	}

	err = db.fDB.Put(GetKeySize(key), encoding.Uint64ToBytes(size+1), nil)
	if err != nil {
		return err
	}

	// push message
	err = db.fDB.Put(
		GetKeyMessage(key, size),
		msg.ToPackage().Bytes(),
		nil,
	)
	if err != nil {
		err := db.fDB.Delete(GetKeyHash(hash), nil)
		if err != nil {
			panic(err)
		}
		err = db.fDB.Delete(GetKeySize(key), nil)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func (db *sKeyValueDB) Load(key []byte, i uint64) local.IMessage {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	data, err := db.fDB.Get(GetKeyMessage(key, i), nil)
	if err != nil {
		return nil
	}
	return local.LoadPackage(data).ToMessage()
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
