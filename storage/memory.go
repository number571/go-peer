package storage

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/encoding"
)

var (
	_ IKeyValueStorage = &sMemoryStorage{}
)

type sMemoryStorage struct {
	fMutex   sync.Mutex
	fMapping map[string][]byte
}

func NewMemoryStorage() IKeyValueStorage {
	return &sMemoryStorage{
		fMapping: make(map[string][]byte),
	}
}

func (store *sMemoryStorage) Set(key, value []byte) error {
	store.fMutex.Lock()
	defer store.fMutex.Unlock()

	store.fMapping[encoding.Base64Encode(key)] = value
	return nil
}

func (store *sMemoryStorage) Get(key []byte) ([]byte, error) {
	store.fMutex.Lock()
	defer store.fMutex.Unlock()

	value, ok := store.fMapping[encoding.Base64Encode(key)]
	if !ok {
		return nil, fmt.Errorf("value undefined")
	}

	return value, nil
}

func (store *sMemoryStorage) Del(key []byte) error {
	store.fMutex.Lock()
	defer store.fMutex.Unlock()

	delete(store.fMapping, encoding.Base64Encode(key))
	return nil
}

func (store *sMemoryStorage) Close() error {
	return nil
}
