package database

import (
	"fmt"
	"os"

	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/local"
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

func (db *sKeyValueDB) Size(key []byte) uint64 {
	data, err := db.fStorage.Get(GetKeySize(key))
	if err != nil {
		return 0
	}

	return encoding.BytesToUint64(data)
}

func (db *sKeyValueDB) Push(key []byte, msg local.IMessage) error {
	// store hash
	hash := msg.Body().Hash()
	_, err := db.fStorage.Get(GetKeyHash(hash))
	if err == nil {
		return fmt.Errorf("hash already exists")
	}

	err = db.fStorage.Set(GetKeyHash(hash), []byte{1})
	if err != nil {
		return err
	}

	// update size
	size := uint64(0)
	bnum, err := db.fStorage.Get(GetKeySize(key))
	if err == nil {
		size = encoding.BytesToUint64(bnum)
	}

	err = db.fStorage.Set(GetKeySize(key), encoding.Uint64ToBytes(size+1))
	if err != nil {
		return err
	}

	// push message
	err = db.fStorage.Set(
		GetKeyMessage(key, size),
		msg.ToPackage().Bytes(),
	)
	if err != nil {
		err := db.fStorage.Del(GetKeyHash(hash))
		if err != nil {
			panic(err)
		}
		err = db.fStorage.Del(GetKeySize(key))
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func (db *sKeyValueDB) Load(key []byte, i uint64) local.IMessage {
	data, err := db.fStorage.Get(GetKeyMessage(key, i))
	if err != nil {
		return nil
	}
	return local.LoadPackage(data).ToMessage()
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
