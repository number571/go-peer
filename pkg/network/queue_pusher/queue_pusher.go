package queue_pusher

import (
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
)

type sQueuePusher struct {
	fMutex sync.Mutex
	fMap   map[string]struct{}
	fQueue []string
	fIndex int
}

func NewQueuePusher(pSettings ISettings) IQueuePusher {
	return &sQueuePusher{
		fQueue: make([]string, pSettings.GetCapacity()),
		fMap:   make(map[string]struct{}, pSettings.GetCapacity()),
	}
}

func (p *sQueuePusher) Push(pKey []byte) bool {
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
	p.fMap[sHash] = struct{}{}

	// increment queue index
	p.fIndex = (p.fIndex + 1) % len(p.fQueue)
	return true
}
