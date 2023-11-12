package storage

import (
	"bytes"
	"encoding/json"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/keybuilder"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/file_system"
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
		encdata, err := file_system.OpenFile(pSettings.GetPath()).Read()
		if err != nil {
			return nil, errors.WrapError(err, "read storage")
		}
		if len(encdata) < cSaltSize {
			return nil, errors.NewError("size of storage < salt size")
		}
		store.fSalt = encdata[:cSaltSize]
	}

	keyBuilder := keybuilder.NewKeyBuilder(pSettings.GetWorkSize(), store.fSalt)
	cipherKey := keyBuilder.Build(pSettings.GetPassword())
	store.fCipher = symmetric.NewAESCipher(cipherKey)

	if !isExist {
		if err := store.Set(nil, nil); err != nil {
			return nil, errors.WrapError(err, "set init storage")
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
			return errors.WrapError(err, "open & decrypt storage")
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
		return nil, errors.NewError("storage undefined")
	}

	// Open and decrypt storage
	mapping, err := p.decrypt()
	if err != nil {
		return nil, errors.WrapError(err, "decrypt storage")
	}

	// Open and decrypt secret
	cipher, mapKey := p.newCipherWithMapKey(pKey)
	encsecret, ok := mapping.FSecrets[mapKey]
	if !ok {
		return nil, errors.NewError("key undefined")
	}
	secret := cipher.DecryptBytes(encsecret)

	return secret, nil
}

func (p *sCryptoStorage) Del(pKey []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	// If storage not exists.
	if !p.isExist() {
		return errors.NewError("storage undefined")
	}

	// Open and decrypt storage
	mapping, err := p.decrypt()
	if err != nil {
		return err
	}

	// Open and decrypt private key
	_, mapKey := p.newCipherWithMapKey(pKey)
	if _, ok := mapping.FSecrets[mapKey]; !ok {
		return errors.NewError("key undefined")
	}

	delete(mapping.FSecrets, mapKey)
	return p.encrypt(mapping)
}

func (p *sCryptoStorage) isExist() bool {
	return file_system.OpenFile(p.fSettings.GetPath()).IsExist()
}

func (p *sCryptoStorage) encrypt(pMapping *storageData) error {
	// Encrypt and save storage
	data, err := json.Marshal(pMapping)
	if err != nil {
		return errors.WrapError(err, "marshal decrypted map")
	}

	err = file_system.OpenFile(p.fSettings.GetPath()).Write(
		bytes.Join(
			[][]byte{p.fSalt, p.fCipher.EncryptBytes(data)},
			[]byte{},
		),
	)
	if err != nil {
		return errors.WrapError(err, "write to storage")
	}

	return nil
}

func (p *sCryptoStorage) decrypt() (*storageData, error) {
	var mapping storageData

	encdata, err := file_system.OpenFile(p.fSettings.GetPath()).Read()
	if err != nil {
		return nil, errors.WrapError(err, "open encrypted storage")
	}

	data := p.fCipher.DecryptBytes(encdata[cSaltSize:])
	err = json.Unmarshal(data, &mapping)
	if err != nil {
		return nil, errors.WrapError(err, "unmarshal decrypt map")
	}

	return &mapping, nil
}

func (p *sCryptoStorage) newCipherWithMapKey(pKey []byte) (symmetric.ICipher, string) {
	keyBuilder := keybuilder.NewKeyBuilder(p.fSettings.GetWorkSize(), p.fSalt)
	cipherKey := keyBuilder.Build(encoding.HexEncode(pKey)) // map key as a password
	return symmetric.NewAESCipher(cipherKey), hashing.NewSHA256Hasher(cipherKey).ToString()
}
