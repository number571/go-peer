package storage

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
)

const (
	cSaltSize = 32
)

var (
	_ IKVStorage = &sCryptoStorage{}
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

func NewCryptoStorage(pSettings ISettings) (IKVStorage, error) {
	store := &sCryptoStorage{
		fSettings: pSettings,
	}

	isExist := store.isExist()

	store.fSalt = random.NewStdPRNG().GetBytes(cSaltSize)
	if isExist {
		encdata, err := os.ReadFile(pSettings.GetPath())
		if err != nil {
			return nil, fmt.Errorf("read storage: %w", err)
		}
		if len(encdata) < cSaltSize {
			return nil, errors.New("size of storage < salt size")
		}
		store.fSalt = encdata[:cSaltSize]
	}

	keyBuilder := keybuilder.NewKeyBuilder(1<<pSettings.GetWorkSize(), store.fSalt)
	cipherKey := keyBuilder.Build(pSettings.GetPassword())
	store.fCipher = symmetric.NewAESCipher(cipherKey)

	if !isExist {
		if err := store.Set(nil, nil); err != nil {
			return nil, fmt.Errorf("set init storage: %w", err)
		}
	}

	return store, nil
}

func (p *sCryptoStorage) GetSettings() ISettings {
	return p.fSettings
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
			return fmt.Errorf("open & decrypt storage: %w", err)
		}
	}

	// Encrypt and save secret into storage
	cipher, mapKey := p.newCipherWithMapKey(pKey)
	mapping.FSecrets[mapKey] = cipher.EncryptBytes(pValue)

	return p.encrypt(mapping)
}

func (p *sCryptoStorage) Get(pKey []byte) ([]byte, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	// If storage not exists.
	if !p.isExist() {
		return nil, errors.New("storage undefined")
	}

	// Open and decrypt storage
	mapping, err := p.decrypt()
	if err != nil {
		return nil, fmt.Errorf("decrypt storage: %w", err)
	}

	// Open and decrypt secret
	cipher, mapKey := p.newCipherWithMapKey(pKey)
	encsecret, ok := mapping.FSecrets[mapKey]
	if !ok {
		return nil, errors.New("key undefined")
	}
	secret := cipher.DecryptBytes(encsecret)

	return secret, nil
}

func (p *sCryptoStorage) Del(pKey []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	// If storage not exists.
	if !p.isExist() {
		return errors.New("storage undefined")
	}

	// Open and decrypt storage
	mapping, err := p.decrypt()
	if err != nil {
		return err
	}

	// Open and decrypt private key
	_, mapKey := p.newCipherWithMapKey(pKey)
	if _, ok := mapping.FSecrets[mapKey]; !ok {
		return errors.New("key undefined")
	}

	delete(mapping.FSecrets, mapKey)
	return p.encrypt(mapping)
}

func (p *sCryptoStorage) isExist() bool {
	_, err := os.Stat(p.fSettings.GetPath())
	return !os.IsNotExist(err)
}

func (p *sCryptoStorage) encrypt(pMapping *storageData) error {
	// Encrypt and save storage
	jsonData := encoding.SerializeJSON(pMapping)
	err := os.WriteFile(
		p.fSettings.GetPath(),
		bytes.Join(
			[][]byte{p.fSalt, p.fCipher.EncryptBytes(jsonData)},
			[]byte{},
		),
		0o644,
	)
	if err != nil {
		return fmt.Errorf("write to storage: %w", err)
	}
	return nil
}

func (p *sCryptoStorage) decrypt() (*storageData, error) {
	var mapping storageData

	encdata, err := os.ReadFile(p.fSettings.GetPath())
	if err != nil {
		return nil, fmt.Errorf("open encrypted storage: %w", err)
	}

	data := p.fCipher.DecryptBytes(encdata[cSaltSize:])
	if err := encoding.DeserializeJSON(data, &mapping); err != nil {
		return nil, fmt.Errorf("unmarshal decrypt map: %w", err)
	}

	return &mapping, nil
}

func (p *sCryptoStorage) newCipherWithMapKey(pKey []byte) (symmetric.ICipher, string) {
	keyBuilder := keybuilder.NewKeyBuilder(1<<p.fSettings.GetWorkSize(), p.fSalt)
	cipherKey := keyBuilder.Build(encoding.HexEncode(pKey)) // map key as a password
	return symmetric.NewAESCipher(cipherKey), hashing.NewSHA256Hasher(cipherKey).ToString()
}
