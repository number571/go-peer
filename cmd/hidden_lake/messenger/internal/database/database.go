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

func NewKeyValueDB(pPath string, pKey []byte) IKeyValueDB {
	db := gp_database.NewLevelDB(
		gp_database.NewSettings(&gp_database.SSettings{
			FPath:      pPath,
			FHashing:   true,
			FCipherKey: pKey,
		}),
	)
	return &sKeyValueDB{
		fDB: &db,
	}
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
		return nil, fmt.Errorf("start > end")
	}

	size := p.getSize(pR)
	if pEnd > size {
		return nil, fmt.Errorf("end > size")
	}

	res := make([]IMessage, 0, pEnd-pStart)
	for i := pStart; i < pEnd; i++ {
		data, err := (*p.fDB).Get(getKeyMessageByEnum(pR, i))
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

func (p *sKeyValueDB) Push(pR IRelation, pMsg IMessage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	_, err := (*p.fDB).Get(getKeyMessageByHash(pR, pMsg.GetSHA256UID()))
	if err == nil {
		return fmt.Errorf("message is already exist")
	}

	err = (*p.fDB).Set(getKeyMessageByHash(pR, pMsg.GetSHA256UID()), []byte{1})
	if err != nil {
		return err
	}

	size := p.getSize(pR)
	err = (*p.fDB).Set(getKeyMessageByEnum(pR, size), pMsg.ToBytes())
	if err != nil {
		return err
	}

	numBytes := encoding.Uint64ToBytes(size + 1)
	return (*p.fDB).Set(getKeySize(pR), numBytes[:])
}

func (p *sKeyValueDB) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return (*p.fDB).Close()
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
