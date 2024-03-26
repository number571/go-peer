package wrapper

import (
	"sync"
)

var (
	_ IWrapper = &sWrapper{}
)

type sWrapper struct {
	fMutex sync.RWMutex
	fValue *interface{}
}

// Used to set values at the time of execution.
func NewWrapper() IWrapper {
	return &sWrapper{fValue: new(interface{})}
}

func (p *sWrapper) Get() interface{} {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return (*p.fValue)
}

func (p *sWrapper) Set(pValue interface{}) IWrapper {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	(*p.fValue) = pValue
	return p
}
