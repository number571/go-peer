package local

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/settings"
)

var (
	_ Storage = &StorageT{}
)

type StorageT struct {
	gs     settings.Settings
	path   string
	salt   []byte
	cipher crypto.Cipher
}

type storageData struct {
	Secrets map[string][]byte `json:"secrets"`
}

func NewStorage(path, password string) Storage {
	store := &StorageT{
		gs:   settings.NewSettings(),
		path: path,
	}

	if store.exists() {
		encdata, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}
		store.salt = encdata[:store.gs.Get(settings.SizeSkey)]
	} else {
		store.salt = crypto.RandBytes(store.gs.Get(settings.SizeSkey))
	}

	ekey := crypto.RaiseEntropy([]byte(password), store.salt, store.gs.Get(settings.SizeWork))
	store.cipher = crypto.NewCipher(ekey)

	if !store.exists() {
		store.Write(nil, "", "")
	}

	return store
}

func (store *StorageT) Write(secret []byte, id, password string) error {
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
		mapping.Secrets = make(map[string][]byte)
	}

	// Encrypt and save private key into storage
	ekey := crypto.RaiseEntropy([]byte(password), bytes.Join(
		[][]byte{
			[]byte(id),
			store.salt,
		},
		[]byte{}), store.gs.Get(settings.SizeWork))
	hash := crypto.NewSHA256(ekey).String()

	cipher := crypto.NewCipher(ekey)
	mapping.Secrets[hash] = cipher.Encrypt(secret)

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

func (store *StorageT) Read(id, password string) ([]byte, error) {
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
	ekey := crypto.RaiseEntropy([]byte(password), bytes.Join(
		[][]byte{
			[]byte(id),
			store.salt,
		},
		[]byte{}), store.gs.Get(settings.SizeWork))
	hash := crypto.NewSHA256(ekey).String()

	encsecret, ok := mapping.Secrets[hash]
	if !ok {
		return nil, fmt.Errorf("error: key undefined")
	}

	cipher := crypto.NewCipher(ekey)
	secret := cipher.Decrypt(encsecret)

	return secret, nil
}

func (store *StorageT) Delete(id, password string) error {
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
	ekey := crypto.RaiseEntropy([]byte(password), bytes.Join(
		[][]byte{
			[]byte(id),
			store.salt,
		},
		[]byte{}), store.gs.Get(settings.SizeWork))
	hash := crypto.NewSHA256(ekey).String()

	_, ok := mapping.Secrets[hash]
	if !ok {
		return fmt.Errorf("error: key undefined")
	}

	delete(mapping.Secrets, hash)

	return nil
}

func (store *StorageT) exists() bool {
	_, err := os.Stat(store.path)
	return !os.IsNotExist(err)
}

func (store *StorageT) decrypt() (storageData, error) {
	var mapping storageData

	encdata, err := ioutil.ReadFile(store.path)
	if err != nil {
		return storageData{}, err
	}

	data := store.cipher.Decrypt(encdata[store.gs.Get(settings.SizeSkey):])
	err = json.Unmarshal(data, &mapping)
	if err != nil {
		return storageData{}, err
	}

	return mapping, nil
}
