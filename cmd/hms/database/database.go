package database

import (
	"fmt"
	"os"
	"sync"

	"github.com/number571/go-peer/crypto/hashing"
	gp_database "github.com/number571/go-peer/database"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/local/message"
)

type sStorage struct {
	fPath string
	fDB   gp_database.IKeyValueDB
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
	dbHashes := gp_database.NewLevelDB(hPath)
	if dbHashes == nil {
		panic("storage (hashes) is nil")
	}
	dbMessages := gp_database.NewLevelDB(mPath)
	if dbMessages == nil {
		panic("storage (messages) is nil")
	}
	return &sKeyValueDB{
		fHashes: &sStorage{
			fPath: hPath,
			fDB:   dbHashes,
		},
		fMessages: &sStorage{
			fPath: mPath,
			fDB:   dbMessages,
		},
	}
}

func (db *sKeyValueDB) Size(key []byte) (uint64, error) {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	if len(key) != hashing.GSHA256Size {
		return 0, fmt.Errorf("key size invalid")
	}

	data, err := db.fMessages.fDB.Get(getKeySize(key))
	if err != nil {
		return 0, nil
	}

	return encoding.BytesToUint64(data), nil
}

func (db *sKeyValueDB) Push(key []byte, msg message.IMessage) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	if len(key) != hashing.GSHA256Size {
		return fmt.Errorf("key size invalid")
	}

	// store hash
	hash := msg.Body().Hash()
	_, err := db.fHashes.fDB.Get(getKeyHash(hash))
	if err == nil {
		return fmt.Errorf("hash already exists")
	}

	err = db.fHashes.fDB.Set(getKeyHash(hash), []byte{1})
	if err != nil {
		return err
	}

	// update size
	size := uint64(0)
	bnum, err := db.fMessages.fDB.Get(getKeySize(key))
	if err == nil {
		size = encoding.BytesToUint64(bnum)
	}

	err = db.fMessages.fDB.Set(getKeySize(key), encoding.Uint64ToBytes(size+1))
	if err != nil {
		return err
	}

	// push message
	err = db.fMessages.fDB.Set(
		getKeyMessage(key, size),
		msg.Bytes(),
	)
	if err != nil {
		err1 := db.fHashes.fDB.Del(getKeyHash(hash))
		if err1 != nil {
			panic(err)
		}
		err2 := db.fMessages.fDB.Del(getKeySize(key))
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

	if len(key) != hashing.GSHA256Size {
		return nil, fmt.Errorf("key size invalid")
	}

	data, err := db.fMessages.fDB.Get(getKeyMessage(key, i))
	if err != nil {
		return nil, fmt.Errorf("message undefined")
	}

	return message.LoadMessage(data), nil
}

func (db *sKeyValueDB) Close() error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	err := db.fMessages.fDB.Close()
	if err != nil {
		return err
	}

	err = db.fHashes.fDB.Close()
	if err != nil {
		return err
	}

	return nil
}

func (db *sKeyValueDB) Clean() error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	err := db.fMessages.fDB.Close()
	if err != nil {
		return err
	}

	err = os.RemoveAll(db.fMessages.fPath)
	if err != nil {
		return err
	}

	db.fMessages.fDB = gp_database.NewLevelDB(db.fMessages.fPath)
	if db.fMessages.fDB == nil {
		panic("storage is nil")
	}

	return nil
}
