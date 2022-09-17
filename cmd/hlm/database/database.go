package database

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/modules/encoding"
	gp_database "github.com/number571/go-peer/modules/storage/database"
)

type sKeyValueDB struct {
	fMutex sync.Mutex
	fDB    gp_database.IKeyValueDB
}

func NewKeyValueDB(path string) IKeyValueDB {
	db := gp_database.NewLevelDB(&gp_database.SSettings{
		FPath: path,
	})
	if db == nil {
		panic("storage (messages) is nil")
	}
	return &sKeyValueDB{
		fDB: db,
	}
}

func (db *sKeyValueDB) Size(rel IRelation) (uint64, error) {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.getSize(rel)
}

func (db *sKeyValueDB) Load(rel IRelation, start, end uint64) ([]string, error) {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	if start > end {
		return nil, fmt.Errorf("start > end")
	}

	size, err := db.getSize(rel)
	if err != nil {
		return nil, err
	}

	if end > size {
		return nil, fmt.Errorf("end > size")
	}

	res := make([]string, 0, end-start+1)
	for i := start; i < end; i++ {
		data, err := db.fDB.Get(getKeyMessage(rel, i))
		if err != nil {
			return nil, fmt.Errorf("message undefined")
		}
		res = append(res, string(data))
	}

	return res, nil
}

func (db *sKeyValueDB) Push(rel IRelation, msg string) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	size, err := db.getSize(rel)
	if err != nil {
		return err
	}

	err = db.fDB.Set(getKeyMessage(rel, size), []byte(msg))
	if err != nil {
		return err
	}

	numBytes := encoding.Uint64ToBytes(size + 1)
	return db.fDB.Set(getKeySize(rel), numBytes[:])
}

func (db *sKeyValueDB) Close() error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fDB.Close()
}

func (db *sKeyValueDB) getSize(rel IRelation) (uint64, error) {
	data, err := db.fDB.Get(getKeySize(rel))
	if err != nil {
		return 0, nil
	}

	res := [encoding.CSizeUint64]byte{}
	copy(res[:], data)
	return encoding.BytesToUint64(res), nil
}
