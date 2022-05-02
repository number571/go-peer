package database

import (
	"fmt"
	"os"

	"github.com/number571/go-peer/crypto"
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

func (db *sKeyValueDB) Size(key []byte) (uint64, error) {
	if len(key) != crypto.HashSize {
		return 0, fmt.Errorf("key size invalid")
	}

	data, err := db.fStorage.Get(getKeySize(key))
	if err != nil {
		return 0, nil
	}

	return encoding.BytesToUint64(data), nil
}

func (db *sKeyValueDB) Push(key []byte, msg local.IMessage) error {
	if len(key) != crypto.HashSize {
		return fmt.Errorf("key size invalid")
	}

	// store hash
	hash := msg.Body().Hash()
	_, err := db.fStorage.Get(getKeyHash(hash))
	if err == nil {
		return fmt.Errorf("hash already exists")
	}

	err = db.fStorage.Set(getKeyHash(hash), []byte{1})
	if err != nil {
		return err
	}

	// update size
	size := uint64(0)
	bnum, err := db.fStorage.Get(getKeySize(key))
	if err == nil {
		size = encoding.BytesToUint64(bnum)
	}

	err = db.fStorage.Set(getKeySize(key), encoding.Uint64ToBytes(size+1))
	if err != nil {
		return err
	}

	// push message
	err = db.fStorage.Set(
		getKeyMessage(key, size),
		msg.ToPackage().Bytes(),
	)
	if err != nil {
		err := db.fStorage.Del(getKeyHash(hash))
		if err != nil {
			panic(err)
		}
		err = db.fStorage.Del(getKeySize(key))
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func (db *sKeyValueDB) Load(key []byte, i uint64) (local.IMessage, error) {
	if len(key) != crypto.HashSize {
		return nil, fmt.Errorf("key size invalid")
	}

	data, err := db.fStorage.Get(getKeyMessage(key, i))
	if err != nil {
		return nil, fmt.Errorf("message undefined")
	}

	return local.LoadPackage(data).ToMessage(), nil
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
