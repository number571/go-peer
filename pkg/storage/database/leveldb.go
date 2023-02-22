package database

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	_ IKeyValueDB = &sLevelDB{}
	_ IIterator   = &sLevelDBIterator{}
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

func NewLevelDB(sett ISettings) IKeyValueDB {
	db, err := leveldb.OpenFile(sett.GetPath(), nil)
	if err != nil {
		return nil
	}
	salt, err := db.Get(sett.GetSaltKey(), nil)
	if err != nil {
		salt = random.NewStdPRNG().Bytes(symmetric.CAESKeySize)
		if err := db.Put(sett.GetSaltKey(), salt, nil); err != nil {
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
	return db.fSettings
}

func (db *sLevelDB) Set(key []byte, value []byte) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fDB.Put(
		db.tryHash(key),
		doEncrypt(db.fCipher, value),
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

	return tryDecrypt(
		db.fCipher,
		encBytes,
	)
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
func (db *sLevelDB) Iter(prefix []byte) IIterator {
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

	decBytes, err := tryDecrypt(
		iter.fCipher,
		iter.fIter.Value(),
	)
	if err != nil {
		return nil
	}
	return decBytes
}

func (iter *sLevelDBIterator) Close() {
	iter.fMutex.Lock()
	defer iter.fMutex.Unlock()

	iter.fIter.Release()
}

func doEncrypt(cipher symmetric.ICipher, dataBytes []byte) []byte {
	return bytes.Join(
		[][]byte{
			hashing.NewHMACSHA256Hasher(
				cipher.Bytes(),
				dataBytes,
			).Bytes(),
			cipher.Encrypt(dataBytes),
		},
		[]byte{},
	)
}

func tryDecrypt(cipher symmetric.ICipher, encBytes []byte) ([]byte, error) {
	if len(encBytes) < hashing.CSHA256Size+symmetric.CAESBlockSize {
		return nil, fmt.Errorf("incorrect size of encrypted data")
	}

	decBytes := cipher.Decrypt(encBytes[hashing.CSHA256Size:])
	if decBytes == nil {
		return nil, fmt.Errorf("failed decrypt message")
	}

	gotHashed := encBytes[:hashing.CSHA256Size]
	newHashed := hashing.NewHMACSHA256Hasher(
		cipher.Bytes(),
		decBytes,
	).Bytes()

	if !bytes.Equal(gotHashed, newHashed) {
		return nil, fmt.Errorf("incorrect hash of decrypted data")
	}

	return decBytes, nil
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
