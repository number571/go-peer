package database

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/number571/go-peer/modules/crypto/hashing"
	"github.com/number571/go-peer/modules/crypto/random"
	"github.com/number571/go-peer/modules/crypto/symmetric"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	cSaltKey = "go-peer/salt"
)

var (
	_ IKeyValueDB = &sLevelDB{}
	_ iIterator   = &sLevelDBIterator{}
)

type sLevelDB struct {
	fMutex    sync.Mutex
	fSalt     []byte
	fDB       *leveldb.DB
	fSettings ISettings
	fCipher   symmetric.ICipher
}

type sLevelDBIterator struct {
	fMutex  sync.Mutex
	fIter   iterator.Iterator
	fCipher symmetric.ICipher
}

func NewLevelDB(sett ISettings, path string) IKeyValueDB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	salt, err := db.Get([]byte(cSaltKey), nil)
	if err != nil {
		salt = random.NewStdPRNG().Bytes(symmetric.CAESKeySize)
		err := db.Put(
			[]byte(cSaltKey),
			salt,
			nil,
		)
		if err != nil {
			return nil
		}
	}
	return &sLevelDB{
		fSalt:     salt,
		fDB:       db,
		fSettings: sett,
		fCipher:   symmetric.NewAESCipher(sett.GetCipherKey()),
	}
}

func (db *sLevelDB) Settings() ISettings {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fSettings
}

func (db *sLevelDB) Set(key []byte, value []byte) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fDB.Put(
		db.tryHash(key),
		db.fCipher.Encrypt(value),
		nil,
	)
}

func (db *sLevelDB) Get(key []byte) ([]byte, error) {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	encBytes, err := db.fDB.Get(db.tryHash(key), nil)
	if err != nil {
		return nil, err
	}

	decBytes := db.fCipher.Decrypt(encBytes)
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

	if db.fSettings.GetHashing() {
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
	if !db.fSettings.GetHashing() {
		return key
	}
	saltWithKey := bytes.Join(
		[][]byte{
			db.fSalt,
			key,
		},
		[]byte{},
	)
	return hashing.NewSHA256Hasher(saltWithKey).Bytes()
}
