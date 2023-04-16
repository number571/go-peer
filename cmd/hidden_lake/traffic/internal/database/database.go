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

func NewKeyValueDB(pSett ISettings) IKeyValueDB {
	levelDB := gp_database.NewLevelDB(
		gp_database.NewSettings(&gp_database.SSettings{
			FPath: pSett.GetPath(),
		}),
	)
	if levelDB == nil {
		panic("storage (hashes) is nil")
	}
	db := &sKeyValueDB{
		fSettings: pSett,
		fDB:       levelDB,
	}
	db.fPointer = db.getPointer()
	return db
}

func (p *sKeyValueDB) Settings() ISettings {
	return p.fSettings
}

func (p *sKeyValueDB) Hashes() ([]string, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	msgsLimit := p.Settings().GetLimitMessages()
	res := make([]string, 0, msgsLimit)
	for i := uint64(0); i < msgsLimit; i++ {
		hash, err := p.fDB.Get(getKeyHash(i))
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

func (p *sKeyValueDB) Push(pMsg message.IMessage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	hash := pMsg.GetBody().GetHash()
	if _, err := p.fDB.Get(getKeyMessage(hash)); err == nil {
		return nil
	}

	params := message.NewSettings(&message.SSettings{
		FWorkSize:    p.Settings().GetWorkSize(),
		FMessageSize: p.Settings().GetMessageSize(),
	})
	if !pMsg.IsValid(params) {
		return fmt.Errorf("invalid push message")
	}

	// delete old message
	keyHash := getKeyHash(p.getPointer())
	if hash, err := p.fDB.Get(keyHash); err == nil {
		keyMsg := getKeyMessage(hash)
		if err := p.fDB.Del(keyMsg); err != nil {
			return err
		}
	}

	// rewrite hash's field
	newHash := pMsg.GetBody().GetHash()
	if err := p.fDB.Set(keyHash, newHash); err != nil {
		return err
	}

	// write message
	keyMsg := getKeyMessage(newHash)
	if err := p.fDB.Set(keyMsg, pMsg.ToBytes()); err != nil {
		return err
	}

	// update pointer
	if err := p.incPointer(); err != nil {
		return err
	}

	return nil
}

func (p *sKeyValueDB) Load(pStrHash string) (message.IMessage, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	hash := encoding.HexDecode(pStrHash)
	if len(hash) != hashing.CSHA256Size {
		return nil, fmt.Errorf("key size invalid")
	}

	data, err := p.fDB.Get(getKeyMessage(hash))
	if err != nil {
		return nil, fmt.Errorf("message undefined")
	}

	msg := message.LoadMessage(
		message.NewSettings(&message.SSettings{
			FWorkSize:    p.Settings().GetWorkSize(),
			FMessageSize: p.Settings().GetMessageSize(),
		}),
		data,
	)
	if msg == nil {
		panic("message is nil")
	}

	return msg, nil
}

func (p *sKeyValueDB) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fDB.Close()
}

func (p *sKeyValueDB) getPointer() uint64 {
	data, err := p.fDB.Get(getKeyPointer())
	if err != nil {
		return 0
	}

	res := [encoding.CSizeUint64]byte{}
	copy(res[:], data)
	return encoding.BytesToUint64(res)
}

func (p *sKeyValueDB) incPointer() error {
	msgsLimit := p.Settings().GetLimitMessages()
	res := encoding.Uint64ToBytes((p.getPointer() + 1) % msgsLimit)
	return p.fDB.Set(getKeyPointer(), res[:])
}
