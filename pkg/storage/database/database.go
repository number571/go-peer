package database

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/storage"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	cSaltKey = "__SALT__"
)

var (
	_ IKVDatabase = &sKeyValueDB{}
)

type sKeyValueDB struct {
	fMutex    sync.Mutex
	fSalt     []byte
	fDB       *leveldb.DB
	fSettings storage.ISettings
	fCipher   symmetric.ICipher
}

func NewKeyValueDB(pSett storage.ISettings) (IKVDatabase, error) {
	db, err := leveldb.OpenFile(pSett.GetPath(), nil)
	if err != nil {
		fmt.Println(pSett.GetPath())
		return nil, errors.WrapError(err, "open database")
	}

	salt, err := db.Get([]byte(cSaltKey), nil)
	if err != nil {
		if !errors.HasError(err, leveldb.ErrNotFound) {
			return nil, errors.WrapError(err, "read salt")
		}
		salt = random.NewStdPRNG().GetBytes(symmetric.CAESKeySize)
		db.Put([]byte(cSaltKey), salt, nil)
	}

	return &sKeyValueDB{
		fSalt:     salt,
		fDB:       db,
		fSettings: pSett,
		fCipher:   symmetric.NewAESCipher(pSett.GetCipherKey()),
	}, nil
}

func (p *sKeyValueDB) GetSettings() storage.ISettings {
	return p.fSettings
}

func (p *sKeyValueDB) Set(pKey []byte, pValue []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if err := p.fDB.Put(p.tryHash(pKey), doEncrypt(p.fCipher, pValue), nil); err != nil {
		return errors.WrapError(err, "insert key/value to database")
	}
	return nil
}

func (p *sKeyValueDB) Get(pKey []byte) ([]byte, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	encValue, err := p.fDB.Get(p.tryHash(pKey), nil)
	if err != nil {
		return nil, errors.WrapError(err, "read value by key")
	}

	if encValue == nil {
		return nil, errors.NewError("undefined value")
	}

	return tryDecrypt(
		p.fCipher,
		encValue,
	)
}

func (p *sKeyValueDB) Del(pKey []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if err := p.fDB.Delete(p.tryHash(pKey), nil); err != nil {
		return errors.WrapError(err, "delete value by key")
	}
	return nil
}

func (p *sKeyValueDB) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if err := p.fDB.Close(); err != nil {
		return errors.WrapError(err, "close database")
	}
	return nil
}

func (p *sKeyValueDB) tryHash(pKey []byte) []byte {
	if !p.fSettings.GetHashing() {
		return pKey
	}
	return hashing.NewHMACSHA256Hasher(p.fSalt, pKey).ToBytes()
}
