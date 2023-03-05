package database

import (
	"fmt"
	"sync"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	gp_database "github.com/number571/go-peer/pkg/storage/database"
)

type sKeyValueDB struct {
	fMutex   sync.Mutex
	fPointer uint64

	fSettings ISettings
	fDB       gp_database.IKeyValueDB
}

func NewKeyValueDB(sett ISettings) IKeyValueDB {
	levelDB := gp_database.NewLevelDB(
		gp_database.NewSettings(&gp_database.SSettings{
			FPath: sett.GetPath(),
		}),
	)
	if levelDB == nil {
		panic("storage (hashes) is nil")
	}
	db := &sKeyValueDB{
		fSettings: sett,
		fDB:       levelDB,
	}
	db.fPointer = db.getPointer()
	return db
}

func (db *sKeyValueDB) Settings() ISettings {
	return db.fSettings
}

func (db *sKeyValueDB) Hashes() ([]string, error) {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	msgsLimit := db.Settings().GetLimitMessages()
	res := make([]string, 0, msgsLimit)
	for i := uint64(0); i < msgsLimit; i++ {
		hash, err := db.fDB.Get(getKeyHash(i))
		if err != nil {
			break
		}
		if len(hash) != hashing.CSHA256Size {
			return nil, fmt.Errorf("incorrect hash size")
		}
		res = append(res, encoding.HexEncode(hash))
	}

	return res, nil
}

func (db *sKeyValueDB) Push(msg message.IMessage) error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	hash := msg.GetBody().GetHash()
	if _, err := db.fDB.Get(getKeyMessage(hash)); err == nil {
		return nil
	}

	params := message.NewParams(
		db.Settings().GetMessageSize(),
		db.Settings().GetWorkSize(),
	)
	if !msg.IsValid(params) {
		return fmt.Errorf("invalid push message")
	}

	// delete old message
	keyHash := getKeyHash(db.getPointer())
	if hash, err := db.fDB.Get(keyHash); err == nil {
		if err := db.fDB.Del(hash); err != nil {
			return err
		}
	}

	// rewrite hash's field
	newHash := msg.GetBody().GetHash()
	if err := db.fDB.Set(keyHash, newHash); err != nil {
		return err
	}

	// write message
	keyMsg := getKeyMessage(newHash)
	if err := db.fDB.Set(keyMsg, msg.ToBytes()); err != nil {
		return err
	}

	// update pointer
	if err := db.incPointer(); err != nil {
		return err
	}

	return nil
}

func (db *sKeyValueDB) Load(strHash string) (message.IMessage, error) {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	hash := encoding.HexDecode(strHash)
	if len(hash) != hashing.CSHA256Size {
		return nil, fmt.Errorf("key size invalid")
	}

	data, err := db.fDB.Get(getKeyMessage(hash))
	if err != nil {
		return nil, fmt.Errorf("message undefined")
	}

	msg := message.LoadMessage(
		data,
		message.NewParams(
			db.Settings().GetMessageSize(),
			db.Settings().GetWorkSize(),
		),
	)
	if msg == nil {
		panic("message is nil")
	}

	return msg, nil
}

func (db *sKeyValueDB) Close() error {
	db.fMutex.Lock()
	defer db.fMutex.Unlock()

	return db.fDB.Close()
}

func (db *sKeyValueDB) getPointer() uint64 {
	data, err := db.fDB.Get(getKeyPointer())
	if err != nil {
		return 0
	}

	res := [encoding.CSizeUint64]byte{}
	copy(res[:], data)
	return encoding.BytesToUint64(res)
}

func (db *sKeyValueDB) incPointer() error {
	msgsLimit := db.Settings().GetLimitMessages()
	res := encoding.Uint64ToBytes((db.getPointer() + 1) % msgsLimit)
	return db.fDB.Set(getKeyPointer(), res[:])
}
