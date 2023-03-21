package database

import (
	"sync"

	"github.com/number571/go-peer/pkg/wrapper"
)

var (
	_ IWrapperDB = &sWrapperDB{}
)

type sWrapperDB struct {
	fMutex   sync.Mutex
	fWrapper wrapper.IWrapper
}

func NewWrapperDB() IWrapperDB {
	return &sWrapperDB{fWrapper: wrapper.NewWrapper()}
}

func (p *sWrapperDB) Get() IKeyValueDB {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	db, ok := p.fWrapper.Get().(IKeyValueDB)
	if !ok {
		return nil
	}

	return db
}

func (p *sWrapperDB) Set(pDB IKeyValueDB) IWrapperDB {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fWrapper.Set(pDB)
	return p
}

func (p *sWrapperDB) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	db, ok := p.fWrapper.Get().(IKeyValueDB)
	if !ok {
		return nil
	}

	p.fWrapper.Set(nil)
	return db.Close()
}
