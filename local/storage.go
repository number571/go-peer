package local

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/number571/gopeer"
	"github.com/number571/gopeer/crypto"
)

type Storage struct {
	path   string
	salt   []byte
	cipher crypto.Cipher
}

type storageData struct {
	Secrets map[string][]byte `json:"secrets"`
}

var (
	workSize = gopeer.Get("POWS_DIFF").(uint) // bits
	saltSize = gopeer.Get("SALT_SIZE").(uint) // bytes
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
		store.salt = encdata[:saltSize]
	} else {
		store.salt = crypto.RandBytes(saltSize)
	}

	ekey := crypto.RaiseEntropy([]byte(password), store.salt, workSize)
	store.cipher = crypto.NewCipher(ekey)

	if !store.Exists() {
		store.Write(nil, "", "")
	}

	return store
}

func (store *Storage) Exists() bool {
	_, err := os.Stat(store.path)
	return !os.IsNotExist(err)
}

func (store *Storage) Write(secret []byte, id, password string) error {
	var (
		mapping storageData
		err     error
	)

	// Open and decrypt storage
	if store.Exists() {
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
		[]byte{}), workSize)
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

func (store *Storage) Read(id, password string) ([]byte, error) {
	// If storage not exists.
	if !store.Exists() {
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
		[]byte{}), workSize)
	hash := crypto.NewSHA256(ekey).String()

	encsecret, ok := mapping.Secrets[hash]
	if !ok {
		return nil, fmt.Errorf("error: key undefined")
	}

	cipher := crypto.NewCipher(ekey)
	secret := cipher.Decrypt(encsecret)

	return secret, nil
}

func (store *Storage) Delete(id, password string) error {
	// If storage not exists.
	if !store.Exists() {
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
		[]byte{}), workSize)
	hash := crypto.NewSHA256(ekey).String()

	_, ok := mapping.Secrets[hash]
	if !ok {
		return fmt.Errorf("error: key undefined")
	}

	delete(mapping.Secrets, hash)

	return nil
}

func (store *Storage) decrypt() (storageData, error) {
	var mapping storageData

	encdata, err := ioutil.ReadFile(store.path)
	if err != nil {
		return storageData{}, err
	}

	data := store.cipher.Decrypt(encdata[saltSize:])
	err = json.Unmarshal(data, &mapping)
	if err != nil {
		return storageData{}, err
	}

	return mapping, nil
}
