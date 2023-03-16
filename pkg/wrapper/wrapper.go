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

func (w *sWrapper) Get() interface{} {
	w.fMutex.Lock()
	defer w.fMutex.Unlock()

	return (*w.fValue)
}

func (w *sWrapper) Set(v interface{}) IWrapper {
	w.fMutex.Lock()
	defer w.fMutex.Unlock()

	(*w.fValue) = v
	return w
}
