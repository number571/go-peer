package database

import (
	"sync"

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

	saltValue, err := db.Get([]byte(cSaltKey), nil)
	if err != nil {
		if !errors.HasError(err, leveldb.ErrNotFound) {
			return nil, errors.WrapError(err, "read salt")
		}
		saltValue = random.NewStdPRNG().GetBytes(2 * cSaltSize)
		db.Put([]byte(cSaltKey), saltValue, nil)
	}

	cipherSalt := saltValue[:cSaltSize]
	cipherKeyBuilder := keybuilder.NewKeyBuilder(pSett.GetWorkSize(), cipherSalt)
	cipherKey := cipherKeyBuilder.Build(pSett.GetPassword())

	authSalt := saltValue[cSaltSize:]
	authKeyBuilder := keybuilder.NewKeyBuilder(pSett.GetWorkSize(), authSalt)
	authKey := authKeyBuilder.Build(pSett.GetPassword())

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

	if encValue == nil {
		return nil, errors.NewError("undefined value")
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
