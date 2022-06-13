package database

import (
	"fmt"
	"os"
	"sync"

	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/local/message"
	"github.com/number571/go-peer/storage"
)

type sStorage struct {
	fPath    string
	fStorage storage.IKeyValueStorage
}

type sKeyValueDB struct {
	fMutex    sync.Mutex
	fHashes   *sStorage
	fMessages *sStorage
}

func NewKeyValueDB(path string) IKeyValueDB {
	var (
		hPath = fmt.Sprintf("%s/hashes", path)
		mPath = fmt.Sprintf("%s/messages", path)
	)
	sHashes := storage.NewLevelDBStorage(hPath)
	if sHashes == nil {
		panic("storage (hashes) is nil")
	}
	sMessages := storage.NewLevelDBStorage(mPath)
	if sMessages == nil {
		panic("storage (messages) is nil")
	}
	return &sKeyValueDB{
		fHashes: &sStorage{
			fPath:    hPath,
			fStorage: sHashes,
		},
		fMessages: &sStorage{
			fPath:    mPath,
			fStorage: sMessages,
		},
	}
}

func (db *sKeyValueDB) Size(key []byte) (uint64, error) {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	if len(key) != hashing.HashSize {
		return 0, fmt.Errorf("key size invalid")
	}

	data, err := db.fMessages.fStorage.Get(getKeySize(key))
	if err != nil {
		return 0, nil
	}

	return encoding.BytesToUint64(data), nil
}

func (db *sKeyValueDB) Push(key []byte, msg message.IMessage) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	if len(key) != hashing.HashSize {
		return fmt.Errorf("key size invalid")
	}

	// store hash
	hash := msg.Body().Hash()
	_, err := db.fHashes.fStorage.Get(getKeyHash(hash))
	if err == nil {
		return fmt.Errorf("hash already exists")
	}

	err = db.fHashes.fStorage.Set(getKeyHash(hash), []byte{1})
	if err != nil {
		return err
	}

	// update size
	size := uint64(0)
	bnum, err := db.fMessages.fStorage.Get(getKeySize(key))
	if err == nil {
		size = encoding.BytesToUint64(bnum)
	}

	err = db.fMessages.fStorage.Set(getKeySize(key), encoding.Uint64ToBytes(size+1))
	if err != nil {
		return err
	}

	// push message
	err = db.fMessages.fStorage.Set(
		getKeyMessage(key, size),
		msg.ToPackage().Bytes(),
	)
	if err != nil {
		err1 := db.fHashes.fStorage.Del(getKeyHash(hash))
		if err1 != nil {
			panic(err)
		}
		err2 := db.fMessages.fStorage.Del(getKeySize(key))
		if err2 != nil {
			panic(err)
		}
		return err
	}

	return nil
}

func (db *sKeyValueDB) Load(key []byte, i uint64) (message.IMessage, error) {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	if len(key) != hashing.HashSize {
		return nil, fmt.Errorf("key size invalid")
	}

	data, err := db.fMessages.fStorage.Get(getKeyMessage(key, i))
	if err != nil {
		return nil, fmt.Errorf("message undefined")
	}

	return message.LoadPackage(data).ToMessage(), nil
}

func (db *sKeyValueDB) Close() error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	err := db.fMessages.fStorage.Close()
	if err != nil {
		return err
	}

	err = db.fHashes.fStorage.Close()
	if err != nil {
		return err
	}

	return nil
}

func (db *sKeyValueDB) Clean() error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	err := db.fMessages.fStorage.Close()
	if err != nil {
		return err
	}

	err = os.RemoveAll(db.fMessages.fPath)
	if err != nil {
		return err
	}

	db.fMessages.fStorage = storage.NewLevelDBStorage(db.fMessages.fPath)
	if db.fMessages.fStorage == nil {
		panic("storage is nil")
	}

	return nil
}
