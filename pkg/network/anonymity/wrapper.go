package anonymity

import (
	"sync"

	"github.com/number571/go-peer/pkg/database"
	"github.com/number571/go-peer/pkg/wrapper"
)

var (
	_ IDBWrapper = &sDBWrapper{}
)

type sDBWrapper struct {
	fMutex   sync.Mutex
	fWrapper wrapper.IWrapper
}

func NewDBWrapper() IDBWrapper {
	return &sDBWrapper{fWrapper: wrapper.NewWrapper()}
}

func (p *sDBWrapper) Get() database.IKVDatabase {
	db, ok := p.fWrapper.Get().(database.IKVDatabase)
	if !ok {
		return nil
	}
	return db
}

func (p *sDBWrapper) Set(pDB database.IKVDatabase) IDBWrapper {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fWrapper.Set(pDB)
	return p
}

func (p *sDBWrapper) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	db, ok := p.fWrapper.Get().(database.IKVDatabase)
	if !ok {
		return nil
	}

	// no need merge errors,
	// database is a third-party package
	p.fWrapper.Set(nil)
	return db.Close()
}
