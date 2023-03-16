package anonymity

import (
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

func (w *sWrapperDB) Get() database.IKeyValueDB {
	db, ok := w.fWrapper.Get().(database.IKeyValueDB)
	if !ok {
		return nil
	}
	return db
}

func (w *sWrapperDB) Set(db database.IKeyValueDB) IWrapperDB {
	w.fMutex.Lock()
	defer w.fMutex.Unlock()

	w.fWrapper.Set(db)
	return w
}

func (w *sWrapperDB) Close() error {
	w.fMutex.Lock()
	defer w.fMutex.Unlock()

	db, ok := w.fWrapper.Get().(database.IKeyValueDB)
	if !ok {
		return nil
	}
	return db.Close()
}
