package selector

import (
	"sync"

	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/crypto/random"
)

var (
	_ ISelector = &sSelector{}
)

type sSelector struct {
	fMutex  sync.Mutex
	fValues []asymmetric.IPubKey
}

func NewSelector(values []asymmetric.IPubKey) ISelector {
	copyPubKeys := make([]asymmetric.IPubKey, len(values))
	copy(copyPubKeys, values)
	return &sSelector{fValues: copyPubKeys}
}

func (s *sSelector) Length() uint64 {
	return uint64(len(s.fValues))
}

func (s *sSelector) Return(n uint64) []asymmetric.IPubKey {
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

	rand := random.NewStdPRNG()
	for i := s.Length() - 1; i > 0; i-- {
		j := int(rand.Uint64() % uint64(i+1))
		s.fValues[i], s.fValues[j] = s.fValues[j], s.fValues[i]
	}

	return s
}
