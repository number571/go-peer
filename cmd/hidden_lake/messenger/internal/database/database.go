package database

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
	gp_database "github.com/number571/go-peer/pkg/storage/database"
)

type sKeyValueDB struct {
	fMutex sync.Mutex
	fDB    *gp_database.IKeyValueDB
}

func NewKeyValueDB(path string, key []byte) IKeyValueDB {
	db := gp_database.NewLevelDB(
		gp_database.NewSettings(&gp_database.SSettings{
			FPath:      path,
			FHashing:   true,
			FCipherKey: key,
		}),
	)
	return &sKeyValueDB{
		fDB: &db,
	}
}

func (db *sKeyValueDB) Size(r IRelation) uint64 {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.getSize(r)
}

func (db *sKeyValueDB) Load(r IRelation, start, end uint64) ([]IMessage, error) {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	if start > end {
		return nil, fmt.Errorf("start > end")
	}

	size := db.getSize(r)
	if end > size {
		return nil, fmt.Errorf("end > size")
	}

	res := make([]IMessage, 0, end-start)
	for i := start; i < end; i++ {
		data, err := (*db.fDB).Get(getKeyMessageByEnum(r, i))
		if err != nil {
			return nil, fmt.Errorf("message undefined")
		}
		msg := LoadMessage(data)
		if msg == nil {
			return nil, fmt.Errorf("message is null")
		}
		res = append(res, msg)
	}

	return res, nil
}

func (db *sKeyValueDB) Push(r IRelation, msg IMessage) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	_, err := (*db.fDB).Get(getKeyMessageByHash(r, msg.GetSHA256UID()))
	if err == nil {
		return fmt.Errorf("message is already exist")
	}

	err = (*db.fDB).Set(getKeyMessageByHash(r, msg.GetSHA256UID()), []byte{1})
	if err != nil {
		return err
	}

	size := db.getSize(r)
	err = (*db.fDB).Set(getKeyMessageByEnum(r, size), msg.Bytes())
	if err != nil {
		return err
	}

	numBytes := encoding.Uint64ToBytes(size + 1)
	return (*db.fDB).Set(getKeySize(r), numBytes[:])
}

func (db *sKeyValueDB) Close() error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return (*db.fDB).Close()
}

func (db *sKeyValueDB) getSize(r IRelation) uint64 {
	data, err := (*db.fDB).Get(getKeySize(r))
	if err != nil {
		return 0
	}

	res := [encoding.CSizeUint64]byte{}
	copy(res[:], data)
	return encoding.BytesToUint64(res)
}
