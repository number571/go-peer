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

func NewLevelDB(pSett ISettings) IKeyValueDB {
	db, err := leveldb.OpenFile(pSett.GetPath(), nil)
	if err != nil {
		return nil
	}
	salt, err := db.Get(pSett.GetSaltKey(), nil)
	if err != nil {
		salt = random.NewStdPRNG().GetBytes(symmetric.CAESKeySize)
		if err := db.Put(pSett.GetSaltKey(), salt, nil); err != nil {
			return nil
		}
	}
	return &sLevelDB{
		fSalt:     salt,
		fDB:       db,
		fSettings: pSett,
		fCipher:   symmetric.NewAESCipher(pSett.GetCipherKey()),
	}
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

// Storage in hashing mode can't iterates
func (p *sLevelDB) GetIterator(pPrefix []byte) IIterator {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fSettings.GetHashing() {
		return nil
	}

	return &sLevelDBIterator{
		fIter:   p.fDB.NewIterator(util.BytesPrefix(pPrefix), nil),
		fCipher: p.fCipher,
	}
}

func (p *sLevelDBIterator) Next() bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fIter.Next()
}

func (p *sLevelDBIterator) GetKey() []byte {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fIter.Key()
}

func (p *sLevelDBIterator) GetValue() []byte {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	decBytes, err := tryDecrypt(
		p.fCipher,
		p.fIter.Value(),
	)
	if err != nil {
		return nil
	}
	return decBytes
}

func (p *sLevelDBIterator) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fIter.Release()
	return nil
}

func doEncrypt(pCipher symmetric.ICipher, pDataBytes []byte) []byte {
	return bytes.Join(
		[][]byte{
			hashing.NewHMACSHA256Hasher(
				pCipher.ToBytes(),
				pDataBytes,
			).ToBytes(),
			pCipher.EncryptBytes(pDataBytes),
		},
		[]byte{},
	)
}

func tryDecrypt(pCipher symmetric.ICipher, pEncBytes []byte) ([]byte, error) {
	if len(pEncBytes) < hashing.CSHA256Size+symmetric.CAESBlockSize {
		return nil, fmt.Errorf("incorrect size of encrypted data")
	}

	decBytes := pCipher.DecryptBytes(pEncBytes[hashing.CSHA256Size:])
	if decBytes == nil {
		return nil, fmt.Errorf("failed decrypt message")
	}

	gotHashed := pEncBytes[:hashing.CSHA256Size]
	newHashed := hashing.NewHMACSHA256Hasher(
		pCipher.ToBytes(),
		decBytes,
	).ToBytes()

	if !bytes.Equal(gotHashed, newHashed) {
		return nil, fmt.Errorf("incorrect hash of decrypted data")
	}

	return decBytes, nil
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
