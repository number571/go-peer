package storage

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/number571/go-peer/crypto/entropy"
	"github.com/number571/go-peer/crypto/hashing"
	"github.com/number571/go-peer/crypto/random"
	"github.com/number571/go-peer/crypto/symmetric"
	"github.com/number571/go-peer/settings"
	"github.com/number571/go-peer/utils"
)

var (
	_ IKeyValueStorage = &sCryptoStorage{}
)

type sCryptoStorage struct {
	fPath     string
	fSalt     []byte
	fWorkSize uint64
	fCipher   symmetric.ICipher
}

type storageData struct {
	FSecrets map[string][]byte `json:"secrets"`
}

// Settings must contain (CSizeSkey, CSizeWork).
func NewCryptoStorage(path string, key []byte, workSize uint64) IKeyValueStorage {
	store := &sCryptoStorage{
		fPath:     path,
		fWorkSize: workSize,
	}

	if store.exists() {
		encdata, err := utils.OpenFile(path).Read()
		if err != nil {
			return nil
		}
		store.fSalt = encdata[:settings.CSizeSymmKey]
	} else {
		store.fSalt = random.NewStdPRNG().Bytes(settings.CSizeSymmKey)
	}

	entropy := entropy.NewEntropy(store.fWorkSize)
	ekey := entropy.Raise(key, store.fSalt)
	store.fCipher = symmetric.NewAESCipher(ekey)

	if !store.exists() {
		store.Set(nil, nil)
	}

	_, err := store.decrypt()
	if err != nil {
		return nil
	}

	return store
}

func (store *sCryptoStorage) Set(key, value []byte) error {
	var (
		mapping storageData
		err     error
	)

	// Open and decrypt storage
	if store.exists() {
		mapping, err = store.decrypt()
		if err != nil {
			return err
		}
	} else {
		mapping.FSecrets = make(map[string][]byte)
	}

	// Encrypt and save private key into storage
	entropy := entropy.NewEntropy(store.fWorkSize)
	ekey := entropy.Raise(key, store.fSalt)
	hash := hashing.NewSHA256Hasher(ekey).String()

	cipher := symmetric.NewAESCipher(ekey)
	mapping.FSecrets[hash] = cipher.Encrypt(value)

	// Encrypt and save storage
	data, err := json.Marshal(&mapping)
	if err != nil {
		return err
	}

	err = utils.OpenFile(store.fPath).Write(
		bytes.Join(
			[][]byte{store.fSalt, store.fCipher.Encrypt(data)},
			[]byte{},
		),
	)
	if err != nil {
		return err
	}

	return nil
}

func (store *sCryptoStorage) Get(key []byte) ([]byte, error) {
	// If storage not exists.
	if !store.exists() {
		return nil, fmt.Errorf("error: storage undefined")
	}

	// Open and decrypt storage
	mapping, err := store.decrypt()
	if err != nil {
		return nil, err
	}

	// Open and decrypt private key
	entropy := entropy.NewEntropy(store.fWorkSize)
	ekey := entropy.Raise(key, store.fSalt)
	hash := hashing.NewSHA256Hasher(ekey).String()

	encsecret, ok := mapping.FSecrets[hash]
	if !ok {
		return nil, fmt.Errorf("error: key undefined")
	}

	cipher := symmetric.NewAESCipher(ekey)
	secret := cipher.Decrypt(encsecret)

	return secret, nil
}

func (store *sCryptoStorage) Del(key []byte) error {
	// If storage not exists.
	if !store.exists() {
		return fmt.Errorf("error: storage undefined")
	}

	// Open and decrypt storage
	mapping, err := store.decrypt()
	if err != nil {
		return err
	}

	// Open and decrypt private key
	entropy := entropy.NewEntropy(store.fWorkSize)
	ekey := entropy.Raise(key, store.fSalt)
	hash := hashing.NewSHA256Hasher(ekey).String()

	_, ok := mapping.FSecrets[hash]
	if !ok {
		return fmt.Errorf("error: key undefined")
	}

	delete(mapping.FSecrets, hash)
	return nil
}

func (store *sCryptoStorage) exists() bool {
	return utils.OpenFile(store.fPath).IsExist()
}

func (store *sCryptoStorage) decrypt() (storageData, error) {
	var mapping storageData

	encdata, err := utils.OpenFile(store.fPath).Read()
	if err != nil {
		return storageData{}, err
	}

	data := store.fCipher.Decrypt(encdata[settings.CSizeSymmKey:])
	err = json.Unmarshal(data, &mapping)
	if err != nil {
		return storageData{}, err
	}

	return mapping, nil
}
