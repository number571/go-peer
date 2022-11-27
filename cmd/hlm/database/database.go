package database

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/encoding"
	gp_database "github.com/number571/go-peer/modules/storage/database"
)

type sKeyValueDB struct {
	fMutex sync.Mutex
	fDB    gp_database.IKeyValueDB
}

func NewKeyValueDB(path, password string) IKeyValueDB {
	db := gp_database.NewLevelDB(
		gp_database.NewSettings(&gp_database.SSettings{
			FPath:      path,
			FHashing:   true,
			FCipherKey: []byte(password),
		}),
	)
	if db == nil {
		return nil
	}
	return &sKeyValueDB{
		fDB: db,
	}
}

func (db *sKeyValueDB) Size(pubKey asymmetric.IPubKey) uint64 {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.getSize(pubKey)
}

func (db *sKeyValueDB) Load(pubKey asymmetric.IPubKey, start, end uint64) ([]IMessage, error) {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	if start > end {
		return nil, fmt.Errorf("start > end")
	}

	size := db.getSize(pubKey)
	if end > size {
		return nil, fmt.Errorf("end > size")
	}

	res := make([]IMessage, 0, end-start)
	for i := start; i < end; i++ {
		data, err := db.fDB.Get(getKeyMessage(pubKey, i))
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

func (db *sKeyValueDB) Push(pubKey asymmetric.IPubKey, msg IMessage) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	size := db.getSize(pubKey)
	err := db.fDB.Set(getKeyMessage(pubKey, size), msg.Bytes())
	if err != nil {
		return err
	}

	numBytes := encoding.Uint64ToBytes(size + 1)
	return db.fDB.Set(getKeySize(pubKey), numBytes[:])
}

func (db *sKeyValueDB) Close() error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fDB.Close()
}

func (db *sKeyValueDB) getSize(pubKey asymmetric.IPubKey) uint64 {
	data, err := db.fDB.Get(getKeySize(pubKey))
	if err != nil {
		return 0
	}

	res := [encoding.CSizeUint64]byte{}
	copy(res[:], data)
	return encoding.BytesToUint64(res)
}
