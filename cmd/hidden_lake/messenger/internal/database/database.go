package database

import (
	"errors"
	"fmt"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/storage"
	gp_database "github.com/number571/go-peer/pkg/storage/database"
)

type sKeyValueDB struct {
	fMutex sync.Mutex
	fDB    *gp_database.IKVDatabase
}

func NewKeyValueDB(pSettings storage.ISettings) (IKVDatabase, error) {
	db, err := gp_database.NewKVDatabase(pSettings)
	if err != nil {
		return nil, fmt.Errorf("new key/value database: %w", err)
	}
	return &sKeyValueDB{
		fDB: &db,
	}, nil
}

func (p *sKeyValueDB) Size(pR IRelation) uint64 {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.getSize(pR)
}

func (p *sKeyValueDB) Load(pR IRelation, pStart, pEnd uint64) ([]IMessage, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if pStart > pEnd {
		return nil, errors.New("start > end")
	}

	size := p.getSize(pR)
	if pEnd > size {
		return nil, errors.New("end > size")
	}

	res := make([]IMessage, 0, pEnd-pStart)
	for i := pStart; i < pEnd; i++ {
		data, err := (*p.fDB).Get(getKeyMessageByEnum(pR, i))
		if err != nil {
			return nil, fmt.Errorf("read message: %w", err)
		}
		msg := LoadMessage(data)
		if msg == nil {
			return nil, errors.New("message is null")
		}
		res = append(res, msg)
	}

	return res, nil
}

func (p *sKeyValueDB) Push(pR IRelation, pMsg IMessage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	size := p.getSize(pR)
	numBytes := encoding.Uint64ToBytes(size + 1)
	if err := (*p.fDB).Set(getKeySize(pR), numBytes[:]); err != nil {
		return fmt.Errorf("set size of message to database: %w", err)
	}

	if err := (*p.fDB).Set(getKeyMessageByEnum(pR, size), pMsg.ToBytes()); err != nil {
		return fmt.Errorf("set message to database: %w", err)
	}

	return nil
}

func (p *sKeyValueDB) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if err := (*p.fDB).Close(); err != nil {
		return fmt.Errorf("close KV database: %w", err)
	}
	return nil
}

func (p *sKeyValueDB) getSize(pR IRelation) uint64 {
	data, err := (*p.fDB).Get(getKeySize(pR))
	if err != nil {
		return 0
	}

	res := [encoding.CSizeUint64]byte{}
	copy(res[:], data)
	return encoding.BytesToUint64(res)
}
