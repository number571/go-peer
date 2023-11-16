package database

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/errors"
	"github.com/number571/go-peer/pkg/queue_set"
)

var (
	_ IKVDatabase = &sVoidKeyValueDB{}
)

type sVoidKeyValueDB struct {
	fSettings ISettings
	fQueueSet queue_set.IQueueSet
}

func NewVoidKeyValueDB(pSett ISettings) IKVDatabase {
	return &sVoidKeyValueDB{
		fSettings: pSett,
		fQueueSet: queue_set.NewQueueSet(
			queue_set.NewSettings(&queue_set.SSettings{
				FCapacity: pSett.GetMessagesCapacity(),
			}),
		),
	}
}

func (p *sVoidKeyValueDB) Settings() ISettings {
	return p.fSettings
}

func (p *sVoidKeyValueDB) Hashes() ([][]byte, error) {
	return nil, nil
}

func (p *sVoidKeyValueDB) Push(pMsg message.IMessage) error {
	if ok := p.fQueueSet.Push(pMsg.GetBody().GetHash(), pMsg.ToBytes()); !ok {
		return errors.OrigError(&SIsExistError{})
	}
	return nil
}

func (p *sVoidKeyValueDB) Load(pHash []byte) (message.IMessage, error) {
	if len(pHash) != hashing.CSHA256Size {
		return nil, errors.NewError("key size invalid")
	}

	msgBytes, ok := p.fQueueSet.Load(pHash)
	if !ok {
		return nil, errors.OrigError(&SIsNotExistError{})
	}

	msg := message.LoadMessage(p.fSettings, msgBytes)
	if msg == nil {
		panic("message is nil")
	}

	return msg, nil
}

func (p *sVoidKeyValueDB) Close() error {
	return nil
}
