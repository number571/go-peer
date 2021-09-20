package local

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/number571/gopeer/crypto"
	"github.com/number571/gopeer/encoding"
)

type Storage struct {
	path   string
	salt   []byte
	cipher crypto.Cipher
}

type storageData struct {
	Keys map[string][]byte `json:"keys"`
}

const (
	SALT_SIZE = 32 // bytes
)

func NewStorage(path, password string) *Storage {
	store := &Storage{
		path: path,
	}

	if store.Exists() {
		encdata, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}
		store.salt = encdata[:SALT_SIZE]
	} else {
		store.salt = crypto.Rand(SALT_SIZE)
	}

	ekey := crypto.RaiseEntropy([]byte(password), store.salt, 20)
	store.cipher = crypto.NewCipher(ekey)

	return store
}

func (store *Storage) Exists() bool {
	_, err := os.Stat(store.path)
	return !os.IsNotExist(err)
}

func (store *Storage) Write(priv crypto.PrivKey, password string) error {
	var mapping storageData
	mapping.Keys = make(map[string][]byte)

	// Open and decrypt storage
	if store.Exists() {
		encdata, err := ioutil.ReadFile(store.path)
		if err != nil {
			return err
		}

		data := store.cipher.Decrypt(encdata[SALT_SIZE:])
		err = json.Unmarshal(data, &mapping)
		if err != nil {
			return err
		}
	}

	// Encrypt and save private key into storage
	ekey := crypto.RaiseEntropy([]byte(password), store.salt, 20)
	hash := encoding.Base64Encode(crypto.SumHash(ekey))

	cipher := crypto.NewCipher(ekey)
	mapping.Keys[hash] = cipher.Encrypt(priv.Bytes())

	// Encrypt and save storage
	data, err := json.Marshal(&mapping)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(store.path, bytes.Join(
		[][]byte{
			store.salt,
			store.cipher.Encrypt(data),
		},
		[]byte{}), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (store *Storage) Read(password string) (crypto.PrivKey, error) {
	var mapping storageData

	// If storage not exists.
	if !store.Exists() {
		return nil, fmt.Errorf("error: storage undefined")
	}

	// Open and decrypt storage
	encdata, err := ioutil.ReadFile(store.path)
	if err != nil {
		return nil, err
	}

	data := store.cipher.Decrypt(encdata[SALT_SIZE:])
	err = json.Unmarshal(data, &mapping)
	if err != nil {
		return nil, err
	}

	// Open and decrypt private key
	ekey := crypto.RaiseEntropy([]byte(password), store.salt, 20)
	hash := encoding.Base64Encode(crypto.SumHash(ekey))

	encpriv, ok := mapping.Keys[hash]
	if !ok {
		return nil, fmt.Errorf("error: key undefined")
	}

	cipher := crypto.NewCipher(ekey)
	priv := crypto.LoadPrivKey(cipher.Decrypt(encpriv))

	if priv == nil {
		return nil, fmt.Errorf("error: key parse")
	}

	return priv, nil
}
