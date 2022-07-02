package netanon

import (
	"sync"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/local/payload"
	"github.com/number571/go-peer/settings"
)

var (
	_ iOnline = &sOnline{}
)

type sOnline struct {
	fMutex   sync.Mutex
	fNode    INode
	fEnabled bool
}

func newOnline(node INode) iOnline {
	return &sOnline{
		fNode: node,
	}
}

func (onl *sOnline) Switch(state bool) {
	onl.fMutex.Lock()
	defer onl.fMutex.Unlock()

	if onl.fEnabled == state {
		return
	}
	onl.fEnabled = state

	switch state {
	case true:
		onl.start()
	case false:
		onl.stop()
	}
}

func (onl *sOnline) Status() bool {
	onl.fMutex.Lock()
	defer onl.fMutex.Unlock()

	return onl.fEnabled
}

func (onl *sOnline) start() {
	go func(node INode) {
		maskPing := node.Client().Settings().Get(settings.CMaskPing)
		node.Handle(
			maskPing,
			func(node INode, sender asymmetric.IPubKey, pl payload.IPayload) []byte {
				return encoding.Uint64ToBytes(maskPing)
			},
		)
	}(onl.fNode)
}

func (onl *sOnline) stop() {
	go func(node INode) {
		maskPing := node.Client().Settings().Get(settings.CMaskPing)
		node.Handle(maskPing, nil)
	}(onl.fNode)
}
