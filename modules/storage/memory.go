package storage

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/modules/encoding"
)

var (
	_ IKeyValueStorage = &sMemoryStorage{}
)

type sMemoryStorage struct {
	fMutex   sync.Mutex
	fMaximum uint64
	fMapping map[string][]byte
}

func NewMemoryStorage(max uint64) IKeyValueStorage {
	return &sMemoryStorage{
		fMaximum: max,
		fMapping: make(map[string][]byte),
	}
}

func (store *sMemoryStorage) Set(key, value []byte) error {
	store.fMutex.Lock()
	defer store.fMutex.Unlock()

	if uint64(len(store.fMapping)) > store.fMaximum {
		for k := range store.fMapping {
			delete(store.fMapping, k)
			break
		}
	}

	store.fMapping[encoding.HexEncode(key)] = value
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
