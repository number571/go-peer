package database

import (
	"bytes"
	"errors"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/utils"
	"go.etcd.io/bbolt"
)

var (
	_ IKVDatabase = &sKVDatabase{}
)

const (
	cBucket  = "_BUCKET_"
	cSaltKey = "__SALT__"
	cRandKey = "__RAND__"
	cHashKey = "__HASH__"
)

const (
	cSaltSize = 32
	cRandSize = 32
)

type sKVDatabase struct {
	fMutex    sync.Mutex
	fDB       *bbolt.DB
	fSettings ISettings
	fCipher   symmetric.ICipher
	fAuthKey  []byte
}

func NewKVDatabase(pSett ISettings) (IKVDatabase, error) {
	db, err := bbolt.Open(pSett.GetPath(), 0600, &bbolt.Options{})
	if err != nil {
		return nil, utils.MergeErrors(ErrOpenDB, err)
	}

	saltValue, initValue, err := getSaltValue(db)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetSalt, err)
	}

	cipherKey, authKey := getCipherAuthKeys(pSett, saltValue)
	randValue, hashRand, err := getRandHashValues(db, initValue, authKey)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetHashRand, err)
	}

	newHashRand := hashing.NewHMACSHA256Hasher(authKey, randValue).ToBytes()
	if !bytes.Equal(hashRand, newHashRand) {
		return nil, ErrInvalidHash
	}

	return &sKVDatabase{
		fDB:       db,
		fSettings: pSett,
		fCipher:   symmetric.NewAESCipher(cipherKey),
		fAuthKey:  authKey,
	}, nil
}

func (p *sKVDatabase) GetSettings() ISettings {
	return p.fSettings
}

func (p *sKVDatabase) Set(pKey []byte, pValue []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	key := hashing.NewHMACSHA256Hasher(p.fAuthKey, pKey).ToBytes()
	val := doEncrypt(p.fCipher, p.fAuthKey, pValue)

	if err := setDB(p.fDB, key, val); err != nil {
		return utils.MergeErrors(ErrSetValueDB, err)
	}
	return nil
}

func setDB(pDB *bbolt.DB, pKey []byte, pValue []byte) error {
	return pDB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(cBucket))
		if err != nil {
			return ErrOpenBucket
		}
		return bucket.Put(pKey, pValue)
	})
}

func (p *sKVDatabase) Get(pKey []byte) ([]byte, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	key := hashing.NewHMACSHA256Hasher(p.fAuthKey, pKey).ToBytes()
	encValue, err := getDB(p.fDB, key)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetValueDB, err)
	}

	return tryDecrypt(
		p.fCipher,
		p.fAuthKey,
		encValue,
	)
}

func getDB(pDB *bbolt.DB, pKey []byte) ([]byte, error) {
	var encValue []byte
	err := pDB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(cBucket))
		if b == nil {
			return ErrOpenBucket
		}
		val := b.Get(pKey)
		if val == nil {
			return ErrGetNotFound
		}
		encValue = make([]byte, len(val))
		copy(encValue, val)
		return nil
	})
	return encValue, err
}

func (p *sKVDatabase) Del(pKey []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	key := hashing.NewHMACSHA256Hasher(p.fAuthKey, pKey).ToBytes()
	if err := delDB(p.fDB, key); err != nil {
		return utils.MergeErrors(ErrDelValueDB, err)
	}
	return nil
}

func delDB(pDB *bbolt.DB, pKey []byte) error {
	return pDB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(cBucket))
		if b == nil {
			return ErrOpenBucket
		}
		return b.Delete(pKey)
	})
}

func (p *sKVDatabase) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if err := p.fDB.Close(); err != nil {
		return utils.MergeErrors(ErrCloseDB, err)
	}
	return nil
}

func getSaltValue(pDB *bbolt.DB) ([]byte, bool, error) {
	initValue := false
	saltValue, err := getDB(pDB, []byte(cSaltKey))
	if err != nil {
		if !errors.Is(err, ErrOpenBucket) && !errors.Is(err, ErrGetNotFound) {
			return nil, false, utils.MergeErrors(ErrReadSalt, err)
		}
		initValue = true
		saltValue = random.NewCSPRNG().GetBytes(cSaltSize)
		if err := setDB(pDB, []byte(cSaltKey), saltValue); err != nil {
			return nil, false, utils.MergeErrors(ErrPushSalt, err)
		}
	}
	return saltValue, initValue, nil
}

func getCipherAuthKeys(pSett ISettings, pSaltValue []byte) ([]byte, []byte) {
	keyBuilder := keybuilder.NewKeyBuilder(1<<pSett.GetWorkSize(), pSaltValue)
	buildKeys := keyBuilder.Build(pSett.GetPassword(), 2*symmetric.CAESKeySize)

	cipherKey := buildKeys[:symmetric.CAESKeySize]
	authKey := buildKeys[symmetric.CAESKeySize:]

	return cipherKey, authKey
}

func getRandHashValues(pDB *bbolt.DB, pInitValue bool, pAuthKey []byte) ([]byte, []byte, error) {
	if pInitValue {
		randValue := random.NewCSPRNG().GetBytes(cRandSize)
		if err := setDB(pDB, []byte(cRandKey), randValue); err != nil {
			return nil, nil, utils.MergeErrors(ErrPushRand, err)
		}
		hashRand := hashing.NewHMACSHA256Hasher(pAuthKey, randValue).ToBytes()
		if err := setDB(pDB, []byte(cHashKey), hashRand); err != nil {
			return nil, nil, utils.MergeErrors(ErrPushHashRand, err)
		}
		return randValue, hashRand, nil
	}

	randValue, err := getDB(pDB, []byte(cRandKey))
	if err != nil {
		return nil, nil, utils.MergeErrors(ErrReadRand, err)
	}
	hashRand, err := getDB(pDB, []byte(cHashKey))
	if err != nil {
		return nil, nil, utils.MergeErrors(ErrReadHashRand, err)
	}
	return randValue, hashRand, nil
}
