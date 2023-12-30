package database

import (
	"bytes"
	"errors"
	"fmt"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/storage"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var (
	_ IKVDatabase = &sKVDatabase{}
)

const (
	cSaltKey  = "__SALT__"
	cHashKey  = "__HASH__"
	cSaltSize = 32
)

type sKVDatabase struct {
	fMutex    sync.Mutex
	fDB       *leveldb.DB
	fSettings storage.ISettings
	fCipher   symmetric.ICipher
	fAuthKey  []byte
}

func NewKVDatabase(pSett storage.ISettings) (IKVDatabase, error) {
	path := pSett.GetPath()
	opt := &opt.Options{
		DisableBlockCache: true,
		Strict:            opt.StrictAll,
	}

	db, err := leveldb.OpenFile(path, opt)
	if err != nil {
		openErr := fmt.Errorf("open database: %w", err)
		db, err = tryRecover(path, opt)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", openErr, err)
		}
	}

	isInitSalt := false
	saltValue, err := db.Get([]byte(cSaltKey), nil)
	if err != nil {
		if !errors.Is(err, leveldb.ErrNotFound) {
			return nil, fmt.Errorf("read salt value: %w", err)
		}
		isInitSalt = true
		saltValue = random.NewStdPRNG().GetBytes(2 * cSaltSize)
		if err := db.Put([]byte(cSaltKey), saltValue, nil); err != nil {
			return nil, fmt.Errorf("put salt value: %w", err)
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
			return nil, fmt.Errorf("put salt hash: %w", err)
		}
	}

	gotSaltHash, err := db.Get([]byte(cHashKey), nil)
	if err != nil {
		return nil, fmt.Errorf("read salt hash: %w", err)
	}

	newSaltHash := hashing.NewHMACSHA256Hasher(authKey, saltValue).ToBytes()
	if !bytes.Equal(gotSaltHash, newSaltHash) {
		return nil, errors.New("incorrect salt hash")
	}

	return &sKVDatabase{
		fDB:       db,
		fSettings: pSett,
		fCipher:   symmetric.NewAESCipher(cipherKey),
		fAuthKey:  authKey,
	}, nil
}

func (p *sKVDatabase) GetSettings() storage.ISettings {
	return p.fSettings
}

func (p *sKVDatabase) Set(pKey []byte, pValue []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	key := hashing.NewHMACSHA256Hasher(p.fAuthKey, pKey).ToBytes()
	if err := p.fDB.Put(key, doEncrypt(p.fCipher, p.fAuthKey, pValue), nil); err != nil {
		return fmt.Errorf("insert key/value to database: %w", err)
	}
	return nil
}

func (p *sKVDatabase) Get(pKey []byte) ([]byte, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	key := hashing.NewHMACSHA256Hasher(p.fAuthKey, pKey).ToBytes()
	encValue, err := p.fDB.Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("read value by key: %w", err)
	}

	return tryDecrypt(
		p.fCipher,
		p.fAuthKey,
		encValue,
	)
}

func (p *sKVDatabase) Del(pKey []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	key := hashing.NewHMACSHA256Hasher(p.fAuthKey, pKey).ToBytes()
	if _, err := p.fDB.Get(key, nil); err != nil {
		return fmt.Errorf("read value by key for delete: %w", err)
	}

	if err := p.fDB.Delete(key, nil); err != nil {
		return fmt.Errorf("delete value by key: %w", err)
	}

	return nil
}

func (p *sKVDatabase) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if err := p.fDB.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}
	return nil
}

func tryRecover(path string, opt *opt.Options) (*leveldb.DB, error) {
	db, err := leveldb.RecoverFile(path, opt)
	if err != nil {
		return nil, fmt.Errorf("recover database: %w", err)
	}
	return db, nil
}
