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

func (w *sWrapperDB) Get() IKeyValueDB {
	w.fMutex.Lock()
	defer w.fMutex.Unlock()

	db, ok := w.fWrapper.Get().(IKeyValueDB)
	if !ok {
		return nil
	}

	return db
}

func (w *sWrapperDB) Set(db IKeyValueDB) IWrapperDB {
	w.fMutex.Lock()
	defer w.fMutex.Unlock()

	w.fWrapper.Set(db)
	return w
}

func (w *sWrapperDB) Close() error {
	w.fMutex.Lock()
	defer w.fMutex.Unlock()

	db, ok := w.fWrapper.Get().(IKeyValueDB)
	if !ok {
		return nil
	}

	w.fWrapper.Set(nil)
	return db.Close()
}
