package local

import (
	"sync"

	"github.com/number571/go-peer/crypto"
)

var (
	_ ISelector = &sSelector{}
)

type sSelector struct {
	fMutex  sync.Mutex
	fValues []crypto.IPubKey
}

func NewSelector(values []crypto.IPubKey) ISelector {
	copyPubKeys := make([]crypto.IPubKey, len(values))
	copy(copyPubKeys, values)
	return &sSelector{fValues: copyPubKeys}
}

func (s *sSelector) Length() uint64 {
	return uint64(len(s.fValues))
}

func (s *sSelector) Return(n uint64) []crypto.IPubKey {
	s.fMutex.Lock()
	defer s.fMutex.Unlock()

	if n > s.Length() {
		n = s.Length()
	}

	return s.fValues[:n]
}

func (s *sSelector) Shuffle() ISelector {
	s.fMutex.Lock()
	defer s.fMutex.Unlock()

	rand := crypto.NewPRNG()
	for i := s.Length() - 1; i > 0; i-- {
		j := int(rand.Uint64() % uint64(i+1))
		s.fValues[i], s.fValues[j] = s.fValues[j], s.fValues[i]
	}

	return s
}
