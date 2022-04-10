package network

import (
	"sync"

	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/local"
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

func (onl *sOnline) Switch(state bool) {
	onl.fMutex.Lock()
	defer onl.fMutex.Unlock()

	if onl.fEnabled == state {
		return
	}

	sett := onl.fNode.Client().Settings()
	patt := encoding.Uint64ToBytes(sett.Get(settings.MaskPing))

	switch state {
	case true:
		onl.fNode.Handle(patt, func(node INode, msg local.IMessage) []byte {
			return patt
		})
	case false:
		onl.fNode.Handle(patt, nil)
	}

	onl.fEnabled = state
}

func (onl *sOnline) Status() bool {
	onl.fMutex.Lock()
	defer onl.fMutex.Unlock()

	return onl.fEnabled
}
