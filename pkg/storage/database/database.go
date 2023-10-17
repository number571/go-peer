package database

import (
	"bytes"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/storage"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const (
	cSaltKey  = "__SALT__"
	cHashKey  = "__HASH__"
	cSaltSize = 32
)

var (
	_ IKVDatabase = &sKeyValueDB{}
)

type sKeyValueDB struct {
	fMutex    sync.Mutex
	fDB       *leveldb.DB
	fSettings storage.ISettings
	fCipher   symmetric.ICipher
	fAuthKey  []byte
}

func NewKeyValueDB(pSett storage.ISettings) (IKVDatabase, error) {
	db, err := leveldb.OpenFile(pSett.GetPath(), &opt.Options{
		DisableBlockCache: true,
	})
	if err != nil {
		return nil, errors.WrapError(err, "open database")
	}

	isInitSalt := false
	saltValue, err := db.Get([]byte(cSaltKey), nil)
	if err != nil {
		if !errors.HasError(err, leveldb.ErrNotFound) {
			return nil, errors.WrapError(err, "read salt value")
		}
		isInitSalt = true
		saltValue = random.NewStdPRNG().GetBytes(2 * cSaltSize)
		if err := db.Put([]byte(cSaltKey), saltValue, nil); err != nil {
			return nil, errors.WrapError(err, "put salt value")
		}
	}

	cipherSalt := saltValue[:cSaltSize]
	cipherKeyBuilder := keybuilder.NewKeyBuilder(pSett.GetWorkSize(), cipherSalt)
	cipherKey := cipherKeyBuilder.Build(pSett.GetPassword())

	authSalt := saltValue[cSaltSize:]
	authKeyBuilder := keybuilder.NewKeyBuilder(pSett.GetWorkSize(), authSalt)
	authKey := authKeyBuilder.Build(pSett.GetPassword())

	if isInitSalt {
		saltHash := hashing.NewHMACSHA256Hasher(authKey, saltValue).ToBytes()
		if err := db.Put([]byte(cHashKey), saltHash, nil); err != nil {
			return nil, errors.WrapError(err, "put salt hash")
		}
	}

	gotSaltHash, err := db.Get([]byte(cHashKey), nil)
	if err != nil {
		return nil, errors.WrapError(err, "read salt hash")
	}

	newSaltHash := hashing.NewHMACSHA256Hasher(authKey, saltValue).ToBytes()
	if !bytes.Equal(gotSaltHash, newSaltHash) {
		return nil, errors.WrapError(err, "incorrect salt hash")
	}

	return &sKeyValueDB{
		fDB:       db,
		fSettings: pSett,
		fCipher:   symmetric.NewAESCipher(cipherKey),
		fAuthKey:  authKey,
	}, nil
}

func (p *sKeyValueDB) GetSettings() storage.ISettings {
	return p.fSettings
}

func (p *sKeyValueDB) Set(pKey []byte, pValue []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if bytes.Equal(pKey, []byte(cSaltKey)) || bytes.Equal(pKey, []byte(cHashKey)) {
		return errors.NewError("key is reserved")
	}

	if err := p.fDB.Put(pKey, doEncrypt(p.fCipher, p.fAuthKey, pValue), nil); err != nil {
		return errors.WrapError(err, "insert key/value to database")
	}
	return nil
}

func (p *sKeyValueDB) Get(pKey []byte) ([]byte, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	encValue, err := p.fDB.Get(pKey, nil)
	if err != nil {
		return nil, errors.WrapError(err, "read value by key")
	}

	return tryDecrypt(
		p.fCipher,
		p.fAuthKey,
		encValue,
	)
}

func (p *sKeyValueDB) Del(pKey []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if bytes.Equal(pKey, []byte(cSaltKey)) || bytes.Equal(pKey, []byte(cHashKey)) {
		return errors.NewError("key is reserved")
	}

	if _, err := p.fDB.Get(pKey, nil); err != nil {
		return errors.WrapError(err, "read value by key")
	}

	if err := p.fDB.Delete(pKey, nil); err != nil {
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
