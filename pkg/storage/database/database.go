package database

import (
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/errors"

	"github.com/akrylysov/pogreb"
)

const (
	cSaltKey = "__SALT__"
)

var (
	_ IKeyValueDB = &sKeyValueDB{}
)

type sKeyValueDB struct {
	fMutex    sync.Mutex
	fSalt     []byte
	fDB       *pogreb.DB
	fSettings ISettings
	fCipher   symmetric.ICipher
}

func NewKeyValueDB(pSett ISettings) (IKeyValueDB, error) {
	db, err := pogreb.Open(pSett.GetPath(), nil)
	if err != nil {
		return nil, errors.WrapError(err, "open database")
	}

	salt, err := db.Get([]byte(cSaltKey))
	if err != nil {
		return nil, errors.WrapError(err, "read salt")
	}

	if salt == nil {
		salt = random.NewStdPRNG().GetBytes(symmetric.CAESKeySize)
		db.Put([]byte(cSaltKey), salt)
	}

	return &sKeyValueDB{
		fSalt:     salt,
		fDB:       db,
		fSettings: pSett,
		fCipher:   symmetric.NewAESCipher(pSett.GetCipherKey()),
	}, nil
}

func (p *sKeyValueDB) GetSettings() ISettings {
	return p.fSettings
}

func (p *sKeyValueDB) Set(pKey []byte, pValue []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if err := p.fDB.Put(p.tryHash(pKey), doEncrypt(p.fCipher, pValue)); err != nil {
		return errors.WrapError(err, "insert key/value to database")
	}
	return nil
}

func (p *sKeyValueDB) Get(pKey []byte) ([]byte, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	encValue, err := p.fDB.Get(p.tryHash(pKey))
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

	if err := p.fDB.Delete(p.tryHash(pKey)); err != nil {
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
