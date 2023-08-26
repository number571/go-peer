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
	"github.com/number571/go-peer/pkg/filesystem"
)

/*
	STORAGE ALGORITHM 'CREATE'

	1. 	S = G(N)
		where
			G - generator pseudo random bytes,
			N - count of bytes for generator,
			S - pseudo random bytes (salt).
	2. 	K = KDF(P, S)
		where
			KDF - key derivation function
			K - encryption key,
			P - password,
	3.	EM = E(K, VM)
		where
			E - encryption function,
			EM - encrypted map/storage,
			VM - void map/storage.
*/

/*
	STORAGE ALGORITHM 'SET'

	1. 	M = D(K, EM)
		where
			D - decryption function
			M - map/storage
	2. 	Km = KDF(Ki, S)
		where
			Km - key map/storage
			Ki - input key
	3. 	Vm = E(Km, Vi)
		where
			Vm - value map/storage
			Vi - input value
	4. 	M = SET(H(Km), Vm)
		where
			H - hash function
			SET - set H(Km),Vm to map/storage
	5. 	EM = E(K, M)
*/

/*
	STORAGE ALGORITHM 'GET'

	1. 	M = D(K, EM)
		where
			M - map/storage
	2. 	Km = KDF(Ki, S)
		where
			Km - key map/storage
			Ki - input key
	3. 	Vm = GET(H(Km))
		where
			GET - get Vm from map/storage by Km
	4. 	Vi = D(Km, Vm)
*/

/*
	STORAGE ALGORITHM 'DEL'

	1. 	M = D(K, EM)
		where
			D - decryption function
			M - map/storage
	2. 	Km = KDF(Ki, S)
		where
			Km - key map/storage
			Ki - input key
	3. 	M = DEL(H(Km))
		where
			DEL - delete Vm from map/storage by Km
	4. 	EM = E(K, M)
*/

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

// pSettings.Hashing always = true
func NewCryptoStorage(pSettings ISettings) (IKVStorage, error) {
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

	keyBuilder := keybuilder.NewKeyBuilder(pSettings.GetWorkSize(), store.fSalt)
	cipherKey := keyBuilder.Build(pSettings.GetPassword())
	store.fCipher = symmetric.NewAESCipher(cipherKey)

	if !isExist {
		store.Set(nil, nil)
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
	return filesystem.OpenFile(p.fSettings.GetPath()).IsExist()
}

func (p *sCryptoStorage) encrypt(pMapping *storageData) error {
	// Encrypt and save storage
	data, err := json.Marshal(pMapping)
	if err != nil {
		return errors.WrapError(err, "marshal decrypted map")
	}

	err = filesystem.OpenFile(p.fSettings.GetPath()).Write(
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

	encdata, err := filesystem.OpenFile(p.fSettings.GetPath()).Read()
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
