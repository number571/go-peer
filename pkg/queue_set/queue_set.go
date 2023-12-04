package queue_set

import (
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IQueueSet = &sQueueSet{}
)

type sQueueSet struct {
	fSettings ISettings
	fMutex    sync.RWMutex
	fMap      map[string][]byte
	fQueue    []string
	fIndex    int
}

func NewQueueSet(pSettings ISettings) IQueueSet {
	return &sQueueSet{
		fSettings: pSettings,
		fQueue:    make([]string, pSettings.GetCapacity()),
		fMap:      make(map[string][]byte, pSettings.GetCapacity()),
	}
}

func (p *sQueueSet) GetSettings() ISettings {
	return p.fSettings
}

func (p *sQueueSet) GetQueueKeys() [][]byte {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	queueSlice := p.fQueue[:p.fIndex]
	keys := make([][]byte, 0, len(queueSlice))
	for _, v := range queueSlice {
		h := encoding.HexDecode(v)
		if h == nil {
			panic("decode hex key param")
		}
		keys = append(keys, h)
	}

	return keys
}

func (p *sQueueSet) Push(pKey, pValue []byte) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	// hash already exists in queue
	sHash := encoding.HexEncode(pKey)
	if _, ok := p.fMap[sHash]; ok {
		return false
	}

	// delete old value in queue
	delete(p.fMap, p.fQueue[p.fIndex])

	// push hash to queue
	p.fQueue[p.fIndex] = sHash
	p.fMap[sHash] = pValue

	// increment queue index
	p.fIndex = (p.fIndex + 1) % len(p.fQueue)
	return true
}

func (p *sQueueSet) Load(pKey []byte) ([]byte, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	val, ok := p.fMap[encoding.HexEncode(pKey)]
	return val, ok
}
