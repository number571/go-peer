package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/entropy"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/filesystem"
)

var (
	_ IKeyValueStorage = &sCryptoStorage{}
)

type sCryptoStorage struct {
	fMutex    sync.Mutex
	fPath     string
	fSalt     []byte
	fWorkSize uint64
	fCipher   symmetric.ICipher
}

type storageData struct {
	FSecrets map[string][]byte `json:"secrets"`
}

func NewCryptoStorage(pPath string, pKey []byte, pWorkSize uint64) (IKeyValueStorage, error) {
	store := &sCryptoStorage{
		fPath:     pPath,
		fWorkSize: pWorkSize,
	}

	if store.exists() {
		encdata, err := filesystem.OpenFile(pPath).Read()
		if err != nil {
			return nil, err
		}
		store.fSalt = encdata[:symmetric.CAESKeySize]
	} else {
		store.fSalt = random.NewStdPRNG().GetBytes(symmetric.CAESKeySize)
	}

	entropy := entropy.NewEntropyBooster(store.fWorkSize, store.fSalt)
	store.fCipher = symmetric.NewAESCipher(entropy.BoostEntropy(pKey))

	if !store.exists() {
		store.Set(nil, nil)
	}

	_, err := store.decrypt()
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (p *sCryptoStorage) Set(pKey, pValue []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	var (
		mapping storageData
		err     error
	)

	// Open and decrypt storage
	if p.exists() {
		mapping, err = p.decrypt()
		if err != nil {
			return err
		}
	} else {
		mapping.FSecrets = make(map[string][]byte)
	}

	// Encrypt and save private key into storage
	entropy := entropy.NewEntropyBooster(p.fWorkSize, p.fSalt)
	ekey := entropy.BoostEntropy(pKey)
	hash := hashing.NewSHA256Hasher(ekey).ToString()

	cipher := symmetric.NewAESCipher(ekey)
	mapping.FSecrets[hash] = cipher.EncryptBytes(pValue)

	return p.encrypt(&mapping)
}

func (p *sCryptoStorage) Get(pKey []byte) ([]byte, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	// If storage not exists.
	if !p.exists() {
		return nil, fmt.Errorf("error: storage undefined")
	}

	// Open and decrypt storage
	mapping, err := p.decrypt()
	if err != nil {
		return nil, err
	}

	// Open and decrypt private key
	entropy := entropy.NewEntropyBooster(p.fWorkSize, p.fSalt)
	ekey := entropy.BoostEntropy(pKey)
	hash := hashing.NewSHA256Hasher(ekey).ToString()

	encsecret, ok := mapping.FSecrets[hash]
	if !ok {
		return nil, fmt.Errorf("error: key undefined")
	}

	cipher := symmetric.NewAESCipher(ekey)
	secret := cipher.DecryptBytes(encsecret)

	return secret, nil
}

func (p *sCryptoStorage) Del(pKey []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	// If storage not exists.
	if !p.exists() {
		return fmt.Errorf("error: storage undefined")
	}

	// Open and decrypt storage
	mapping, err := p.decrypt()
	if err != nil {
		return err
	}

	// Open and decrypt private key
	entropy := entropy.NewEntropyBooster(p.fWorkSize, p.fSalt)
	hash := hashing.NewSHA256Hasher(entropy.BoostEntropy(pKey)).ToString()

	_, ok := mapping.FSecrets[hash]
	if !ok {
		return fmt.Errorf("error: key undefined")
	}

	delete(mapping.FSecrets, hash)
	return p.encrypt(&mapping)
}

func (p *sCryptoStorage) exists() bool {
	return filesystem.OpenFile(p.fPath).IsExist()
}

func (p *sCryptoStorage) encrypt(pMapping *storageData) error {
	// Encrypt and save storage
	data, err := json.Marshal(pMapping)
	if err != nil {
		return err
	}

	err = filesystem.OpenFile(p.fPath).Write(
		bytes.Join(
			[][]byte{p.fSalt, p.fCipher.EncryptBytes(data)},
			[]byte{},
		),
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *sCryptoStorage) decrypt() (storageData, error) {
	var mapping storageData

	encdata, err := filesystem.OpenFile(p.fPath).Read()
	if err != nil {
		return storageData{}, err
	}

	data := p.fCipher.DecryptBytes(encdata[symmetric.CAESKeySize:])
	err = json.Unmarshal(data, &mapping)
	if err != nil {
		return storageData{}, err
	}

	return mapping, nil
}
