package local

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/number571/go-peer/cmd/hls/utils"
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/settings"
)

var (
	_ IStorage = &sStorage{}
)

type sStorage struct {
	fSettings settings.ISettings
	fPath     string
	fSalt     []byte
	fCipher   crypto.ICipher
}

type storageData struct {
	FSecrets map[string][]byte `json:"secrets"`
}

func NewStorage(sett settings.ISettings, path string, pasw string) IStorage {
	store := &sStorage{
		fSettings: sett,
		fPath:     path,
	}

	if store.exists() {
		encdata, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}
		store.fSalt = encdata[:store.fSettings.Get(settings.SizeSkey)]
	} else {
		store.fSalt = crypto.NewPRNG().Bytes(store.fSettings.Get(settings.SizeSkey))
	}

	ekey := crypto.RaiseEntropy([]byte(pasw), store.fSalt, store.fSettings.Get(settings.SizeWork))
	store.fCipher = crypto.NewCipher(ekey)

	if !store.exists() {
		store.Write("", "", nil)
	}

	_, err := store.decrypt()
	if err != nil {
		return nil
	}

	return store
}

func (store *sStorage) Write(id string, pasw string, secret []byte) error {
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
	ekey := crypto.RaiseEntropy([]byte(pasw), bytes.Join(
		[][]byte{
			[]byte(id),
			store.fSalt,
		},
		[]byte{}), store.fSettings.Get(settings.SizeWork))
	hash := crypto.NewHasher(ekey).String()

	cipher := crypto.NewCipher(ekey)
	mapping.FSecrets[hash] = cipher.Encrypt(secret)

	// Encrypt and save storage
	data, err := json.Marshal(&mapping)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(store.fPath, bytes.Join(
		[][]byte{
			store.fSalt,
			store.fCipher.Encrypt(data),
		},
		[]byte{}), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (store *sStorage) Read(id string, pasw string) ([]byte, error) {
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
	ekey := crypto.RaiseEntropy([]byte(pasw), bytes.Join(
		[][]byte{
			[]byte(id),
			store.fSalt,
		},
		[]byte{}), store.fSettings.Get(settings.SizeWork))
	hash := crypto.NewHasher(ekey).String()

	encsecret, ok := mapping.FSecrets[hash]
	if !ok {
		return nil, fmt.Errorf("error: key undefined")
	}

	cipher := crypto.NewCipher(ekey)
	secret := cipher.Decrypt(encsecret)

	return secret, nil
}

func (store *sStorage) Delete(id string, pasw string) error {
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
	ekey := crypto.RaiseEntropy([]byte(pasw), bytes.Join(
		[][]byte{
			[]byte(id),
			store.fSalt,
		},
		[]byte{}), store.fSettings.Get(settings.SizeWork))
	hash := crypto.NewHasher(ekey).String()

	_, ok := mapping.FSecrets[hash]
	if !ok {
		return fmt.Errorf("error: key undefined")
	}

	delete(mapping.FSecrets, hash)

	return nil
}

func (store *sStorage) exists() bool {
	return utils.FileIsExist(store.fPath)
}

func (store *sStorage) decrypt() (storageData, error) {
	var mapping storageData

	encdata, err := ioutil.ReadFile(store.fPath)
	if err != nil {
		return storageData{}, err
	}

	data := store.fCipher.Decrypt(encdata[store.fSettings.Get(settings.SizeSkey):])
	err = json.Unmarshal(data, &mapping)
	if err != nil {
		return storageData{}, err
	}

	return mapping, nil
}
