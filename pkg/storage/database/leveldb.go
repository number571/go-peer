package database

import (
	"bytes"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	_ IKeyValueDB = &sLevelDB{}
)

type sLevelDB struct {
	fMutex    sync.Mutex
	fSalt     []byte
	fDB       *leveldb.DB
	fSettings ISettings
	fCipher   symmetric.ICipher
}

func NewLevelDB(pSett ISettings) (IKeyValueDB, error) {
	db, err := leveldb.OpenFile(pSett.GetPath(), nil)
	if err != nil {
		return nil, err
	}
	salt, err := db.Get(pSett.GetSaltKey(), nil)
	if err != nil {
		salt = random.NewStdPRNG().GetBytes(symmetric.CAESKeySize)
		if err := db.Put(pSett.GetSaltKey(), salt, nil); err != nil {
			return nil, err
		}
	}
	return &sLevelDB{
		fSalt:     salt,
		fDB:       db,
		fSettings: pSett,
		fCipher:   symmetric.NewAESCipher(pSett.GetCipherKey()),
	}, nil
}

func (p *sLevelDB) GetSettings() ISettings {
	return p.fSettings
}

func (p *sLevelDB) Set(pKey []byte, pValue []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fDB.Put(
		p.tryHash(pKey),
		doEncrypt(p.fCipher, pValue),
		nil,
	)
}

func (p *sLevelDB) Get(pKey []byte) ([]byte, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	encBytes, err := p.fDB.Get(p.tryHash(pKey), nil)
	if err != nil {
		return nil, err
	}

	return tryDecrypt(
		p.fCipher,
		encBytes,
	)
}

func (p *sLevelDB) Del(pKey []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fDB.Delete(p.tryHash(pKey), nil)
}

func (p *sLevelDB) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fDB.Close()
}

func (p *sLevelDB) tryHash(pKey []byte) []byte {
	if !p.fSettings.GetHashing() {
		return pKey
	}
	saltWithKey := bytes.Join(
		[][]byte{
			p.fSalt,
			pKey,
		},
		[]byte{},
	)
	return hashing.NewSHA256Hasher(saltWithKey).ToBytes()
}
