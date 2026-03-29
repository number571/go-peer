package cache

import (
	"sync"
)

var (
	_ ILRUCache = &sLRUCache{}
)

type sLRUCache struct {
	fMutex sync.RWMutex
	fMap   map[string]interface{}
	fQueue []string
	fIndex uint64
}

func NewLRUCache(pCapacity uint64) ILRUCache {
	return &sLRUCache{
		fQueue: make([]string, pCapacity),
		fMap:   make(map[string]interface{}, pCapacity),
	}
}

func (p *sLRUCache) GetIndex() uint64 {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return p.fIndex
}

func (p *sLRUCache) GetKey(i uint64) (string, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	if uint64(len(p.fQueue)) <= i {
		return "", false
	}

	key := p.fQueue[i]
	return key, len(key) != 0
}

func (p *sLRUCache) Get(pKey string) (interface{}, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	val, ok := p.fMap[pKey]
	return val, ok
}

func (p *sLRUCache) Set(pKey string, pValue interface{}) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	// hash already exists in queue
	if _, ok := p.fMap[pKey]; ok {
		return false
	}

	// delete old value in queue
	delete(p.fMap, p.fQueue[p.fIndex])

	// push hash to queue
	p.fQueue[p.fIndex] = pKey
	p.fMap[pKey] = pValue

	// increment queue index
	p.fIndex = (p.fIndex + 1) % uint64(len(p.fQueue))
	return true
}
