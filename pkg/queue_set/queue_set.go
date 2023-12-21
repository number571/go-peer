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
	fMutex    sync.Mutex
	fMap      map[string][]byte
	fQueue    []string
	fIndex    uint64
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

func (p *sQueueSet) GetIndex() uint64 {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fIndex
}

func (p *sQueueSet) GetKey(i uint64) ([]byte, bool) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	queueLen := uint64(len(p.fQueue))
	if queueLen <= i {
		return nil, false
	}

	hash := encoding.HexDecode(p.fQueue[i])
	if hash == nil {
		return nil, false
	}

	return hash, true
}

func (p *sQueueSet) Push(pKey, pValue []byte) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	// hash already exists in queue
	encKey := encoding.HexEncode(pKey)
	if _, ok := p.fMap[encKey]; ok {
		return false
	}

	// delete old value in queue
	delete(p.fMap, p.fQueue[p.fIndex])

	// push hash to queue
	p.fQueue[p.fIndex] = encKey
	p.fMap[encKey] = pValue

	// increment queue index
	p.fIndex = (p.fIndex + 1) % uint64(len(p.fQueue))
	return true
}

func (p *sQueueSet) Load(pKey []byte) ([]byte, bool) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	val, ok := p.fMap[encoding.HexEncode(pKey)]
	return val, ok
}
