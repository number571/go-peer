package anonymity

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/pkg/storage/database"
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

func (p *sWrapperDB) Get() database.IKVDatabase {
	db, ok := p.fWrapper.Get().(database.IKVDatabase)
	if !ok {
		return nil
	}
	return db
}

func (p *sWrapperDB) Set(pDB database.IKVDatabase) IWrapperDB {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fWrapper.Set(pDB)
	return p
}

func (p *sWrapperDB) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	db, ok := p.fWrapper.Get().(database.IKVDatabase)
	if !ok {
		return nil
	}

	p.fWrapper.Set(nil)
	if err := db.Close(); err != nil {
		return fmt.Errorf("close wrapped database: %w", err)
	}
	return nil
}
