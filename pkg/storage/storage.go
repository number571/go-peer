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

const (
	cSaltSize = symmetric.CAESKeySize
)

var (
	_ IKeyValueStorage = &sCryptoStorage{}
)

type sCryptoStorage struct {
	fMutex    sync.Mutex
	fSalt     []byte
	fSettings ISettings
	fCipher   symmetric.ICipher
}

type storageData struct {
	FSecrets map[string][]byte `json:"secrets"`
}

func NewCryptoStorage(pSettings ISettings) (IKeyValueStorage, error) {
	store := &sCryptoStorage{
		fSettings: pSettings,
	}
	isExist := store.isExist()

	store.fSalt = random.NewStdPRNG().GetBytes(cSaltSize)
	if isExist {
		encdata, err := filesystem.OpenFile(pSettings.GetPath()).Read()
		if err != nil {
			return nil, err
		}
		store.fSalt = encdata[:cSaltSize]
	}

	entropy := entropy.NewEntropyBooster(pSettings.GetWorkSize(), store.fSalt)
	store.fCipher = symmetric.NewAESCipher(entropy.BoostEntropy(pSettings.GetCipherKey()))

	if !isExist {
		store.Set(nil, nil)
	}

	return store, nil
}

func (p *sCryptoStorage) Set(pKey, pValue []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	mapping := new(storageData)
	mapping.FSecrets = make(map[string][]byte)

	// Open and decrypt storage
	if p.isExist() {
		var err error
		mapping, err = p.decrypt()
		if err != nil {
			return err
		}
	}

	// Encrypt and save secret into storage
	cipher, hash := p.newCipherWithKeyHash(pKey)
	mapping.FSecrets[hash] = cipher.EncryptBytes(pValue)

	return p.encrypt(mapping)
}

func (p *sCryptoStorage) Get(pKey []byte) ([]byte, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	// If storage not exists.
	if !p.isExist() {
		return nil, fmt.Errorf("error: storage undefined")
	}

	// Open and decrypt storage
	mapping, err := p.decrypt()
	if err != nil {
		return nil, err
	}

	// Open and decrypt secret
	cipher, hash := p.newCipherWithKeyHash(pKey)
	encsecret, ok := mapping.FSecrets[hash]
	if !ok {
		return nil, fmt.Errorf("error: key undefined")
	}
	secret := cipher.DecryptBytes(encsecret)

	return secret, nil
}

func (p *sCryptoStorage) Del(pKey []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	// If storage not exists.
	if !p.isExist() {
		return fmt.Errorf("error: storage undefined")
	}

	// Open and decrypt storage
	mapping, err := p.decrypt()
	if err != nil {
		return err
	}

	// Open and decrypt private key
	_, hash := p.newCipherWithKeyHash(pKey)
	_, ok := mapping.FSecrets[hash]
	if !ok {
		return fmt.Errorf("error: key undefined")
	}

	delete(mapping.FSecrets, hash)
	return p.encrypt(mapping)
}

func (p *sCryptoStorage) isExist() bool {
	return filesystem.OpenFile(p.fSettings.GetPath()).IsExist()
}

func (p *sCryptoStorage) encrypt(pMapping *storageData) error {
	// Encrypt and save storage
	data, err := json.Marshal(pMapping)
	if err != nil {
		return err
	}

	err = filesystem.OpenFile(p.fSettings.GetPath()).Write(
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

func (p *sCryptoStorage) decrypt() (*storageData, error) {
	var mapping storageData

	encdata, err := filesystem.OpenFile(p.fSettings.GetPath()).Read()
	if err != nil {
		return nil, err
	}

	data := p.fCipher.DecryptBytes(encdata[cSaltSize:])
	err = json.Unmarshal(data, &mapping)
	if err != nil {
		return nil, err
	}

	return &mapping, nil
}

func (p *sCryptoStorage) newCipherWithKeyHash(pKey []byte) (symmetric.ICipher, string) {
	entropy := entropy.NewEntropyBooster(p.fSettings.GetWorkSize(), p.fSalt)
	ekey := entropy.BoostEntropy(pKey)
	hash := hashing.NewSHA256Hasher(ekey).ToString()
	return symmetric.NewAESCipher(ekey), hash
}
