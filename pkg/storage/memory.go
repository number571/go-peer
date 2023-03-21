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

func NewMemoryStorage(pMaximum uint64) IKeyValueStorage {
	return &sMemoryStorage{
		fMaximum:  pMaximum,
		fKeyQueue: make([]string, 0, pMaximum),
		fMapping:  make(map[string][]byte, pMaximum),
	}
}

func (p *sMemoryStorage) Set(pKey, pValue []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if uint64(len(p.fMapping)) >= p.fMaximum {
		delete(p.fMapping, p.fKeyQueue[0])
		p.fKeyQueue = p.fKeyQueue[1:]
	}

	newKey := encoding.HexEncode(pKey)

	p.fKeyQueue = append(p.fKeyQueue, newKey)
	p.fMapping[newKey] = pValue
	return nil
}

func (p *sMemoryStorage) Get(pKey []byte) ([]byte, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	value, ok := p.fMapping[encoding.HexEncode(pKey)]
	if !ok {
		return nil, fmt.Errorf("undefined value by key")
	}

	return value, nil
}

func (p *sMemoryStorage) Del(pKey []byte) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	_, ok := p.fMapping[encoding.HexEncode(pKey)]
	if !ok {
		return fmt.Errorf("undefined value by key")
	}

	delete(p.fMapping, encoding.HexEncode(pKey))
	return nil
}
