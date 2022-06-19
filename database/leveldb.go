package database

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/crypto/symmetric"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	_ IKeyValueDB = &sLevelDB{}
	_ iIterator   = &sLevelDBIterator{}
)

type sLevelDB struct {
	fMutex   sync.Mutex
	fDB      *leveldb.DB
	fCipher  symmetric.ICipher
	fHashing bool
}

type sLevelDBIterator struct {
	fMutex  sync.Mutex
	fIter   iterator.Iterator
	fCipher symmetric.ICipher
}

func NewLevelDB(path string) IKeyValueDB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil
	}
	return &sLevelDB{
		fDB: db,
	}
}

func (db *sLevelDB) WithEncryption(key []byte) IKeyValueDB {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	var cipher symmetric.ICipher
	if key != nil {
		cipher = symmetric.NewAESCipher(key)
	}
	db.fCipher = cipher
	return db
}

func (db *sLevelDB) WithHashing(state bool) IKeyValueDB {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	db.fHashing = state
	return db
}

func (db *sLevelDB) Set(key []byte, value []byte) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fDB.Put(db.tryHash(key), db.tryEncrypt(value), nil)
}

func (db *sLevelDB) Get(key []byte) ([]byte, error) {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	encBytes, err := db.fDB.Get(db.tryHash(key), nil)
	if err != nil {
		return nil, err
	}

	decBytes := db.tryDecrypt(encBytes)
	if decBytes == nil {
		return nil, fmt.Errorf("failed decrypt message")
	}

	return decBytes, nil
}

func (db *sLevelDB) Del(key []byte) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fDB.Delete(db.tryHash(key), nil)
}

func (db *sLevelDB) Close() error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fDB.Close()
}

// Storage in hashing mode can't iterates
func (db *sLevelDB) Iter(prefix []byte) iIterator {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	if db.fHashing {
		return nil
	}

	return &sLevelDBIterator{
		fIter:   db.fDB.NewIterator(util.BytesPrefix(prefix), nil),
		fCipher: db.fCipher,
	}
}

func (iter *sLevelDBIterator) Next() bool {
	iter.fMutex.Lock()
	defer iter.fMutex.Unlock()

	return iter.fIter.Next()
}

func (iter *sLevelDBIterator) Key() []byte {
	iter.fMutex.Lock()
	defer iter.fMutex.Unlock()

	return iter.fIter.Key()
}

func (iter *sLevelDBIterator) Value() []byte {
	iter.fMutex.Lock()
	defer iter.fMutex.Unlock()

	if iter.fCipher == nil {
		return iter.fIter.Value()
	}

	return iter.fCipher.Decrypt(iter.fIter.Value())
}

func (iter *sLevelDBIterator) Close() {
	iter.fMutex.Lock()
	defer iter.fMutex.Unlock()

	iter.fIter.Release()
}

func (db *sLevelDB) tryHash(key []byte) []byte {
	if db.fHashing {
		return hashing.NewSHA256Hasher(key).Bytes()
	}
	return key
}

func (db *sLevelDB) tryEncrypt(value []byte) []byte {
	if db.fCipher != nil {
		return db.fCipher.Encrypt(value)
	}
	return value
}

func (db *sLevelDB) tryDecrypt(value []byte) []byte {
	if db.fCipher != nil {
		return db.fCipher.Decrypt(value)
	}
	return value
}
