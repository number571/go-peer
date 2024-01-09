package anonymity

import (
	"sync"

	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/utils"
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

	p.fWrapper.Set(nil)
	if err := db.Close(); err != nil {
		return utils.MergeErrors(ErrCloseWrapperDB, err)
	}
	return nil
}
