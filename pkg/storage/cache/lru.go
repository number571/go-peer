package cache

import (
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ ILRUCache = &sLRUCache{}
)

type sLRUCache struct {
	fMutex sync.RWMutex
	fMap   map[string][]byte
	fQueue []string
	fIndex uint64
}

func NewLRUCache(pCapacity uint64) ILRUCache {
	return &sLRUCache{
		fQueue: make([]string, pCapacity),
		fMap:   make(map[string][]byte, pCapacity),
	}
}

func (p *sLRUCache) GetIndex() uint64 {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return p.fIndex
}

func (p *sLRUCache) GetKey(i uint64) ([]byte, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	if uint64(len(p.fQueue)) <= i {
		return nil, false
	}

	hash := encoding.HexDecode(p.fQueue[i])
	return hash, len(hash) != 0
}

func (p *sLRUCache) Get(pKey []byte) ([]byte, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	val, ok := p.fMap[encoding.HexEncode(pKey)]
	return val, ok
}

func (p *sLRUCache) Set(pKey, pValue []byte) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	// hash already exists in queue
	hexKey := encoding.HexEncode(pKey)
	if _, ok := p.fMap[hexKey]; ok {
		return false
	}

	// delete old value in queue
	delete(p.fMap, p.fQueue[p.fIndex])

	// push hash to queue
	p.fQueue[p.fIndex] = hexKey
	p.fMap[hexKey] = pValue

	// increment queue index
	p.fIndex = (p.fIndex + 1) % uint64(len(p.fQueue))
	return true
}
