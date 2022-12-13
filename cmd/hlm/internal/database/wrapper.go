package database

import (
	"fmt"
	"sync"
)

var (
	_ IWrapperDB = &sWrapperDB{}
)

type sWrapperDB struct {
	fMutex    sync.Mutex
	fDatabase *IKeyValueDB
}

func NewWrapperDB() IWrapperDB {
	return &sWrapperDB{
		fDatabase: new(IKeyValueDB),
	}
}

func (w *sWrapperDB) Get() IKeyValueDB {
	return (*w.fDatabase)
}

func (w *sWrapperDB) Update(db IKeyValueDB) error {
	w.fMutex.Lock()
	defer w.fMutex.Unlock()

	if (*w.fDatabase) != nil {
		return fmt.Errorf("failed: pointer already exist")
	}

	(*w.fDatabase) = db
	return nil
}

func (w *sWrapperDB) Close() error {
	w.fMutex.Lock()
	defer w.fMutex.Unlock()

	var err error
	if (*w.fDatabase) != nil {
		err = (*w.fDatabase).Close()
	}

	(*w.fDatabase) = nil
	return err
}
