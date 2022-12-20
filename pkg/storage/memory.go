package storage

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IKeyValueStorage = &sMemoryStorage{}
)

type sMemoryStorage struct {
	fMutex    sync.Mutex
	fMaximum  uint64
	fKeyQueue []string
	fMapping  map[string][]byte
}

func NewMemoryStorage(max uint64) IKeyValueStorage {
	return &sMemoryStorage{
		fMaximum:  max,
		fKeyQueue: make([]string, 0, max),
		fMapping:  make(map[string][]byte, max),
	}
}

func (store *sMemoryStorage) Set(key, value []byte) error {
	store.fMutex.Lock()
	defer store.fMutex.Unlock()

	if uint64(len(store.fMapping)) >= store.fMaximum {
		delete(store.fMapping, store.fKeyQueue[0])
		store.fKeyQueue = store.fKeyQueue[1:]
	}

	newKey := encoding.HexEncode(key)

	store.fKeyQueue = append(store.fKeyQueue, newKey)
	store.fMapping[newKey] = value
	return nil
}

func (store *sMemoryStorage) Get(key []byte) ([]byte, error) {
	store.fMutex.Lock()
	defer store.fMutex.Unlock()

	value, ok := store.fMapping[encoding.HexEncode(key)]
	if !ok {
		return nil, fmt.Errorf("undefined value by key")
	}

	return value, nil
}

func (store *sMemoryStorage) Del(key []byte) error {
	store.fMutex.Lock()
	defer store.fMutex.Unlock()

	_, ok := store.fMapping[encoding.HexEncode(key)]
	if !ok {
		return fmt.Errorf("undefined value by key")
	}

	delete(store.fMapping, encoding.HexEncode(key))
	return nil
}
