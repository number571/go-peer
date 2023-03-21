package wrapper

import (
	"sync"
)

var (
	_ IWrapper = &sWrapper{}
)

type sWrapper struct {
	fMutex sync.Mutex
	fValue *interface{}
}

func NewWrapper() IWrapper {
	return &sWrapper{fValue: new(interface{})}
}

func (p *sWrapper) Get() interface{} {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return (*p.fValue)
}

func (p *sWrapper) Set(pValue interface{}) IWrapper {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	(*p.fValue) = pValue
	return p
}
