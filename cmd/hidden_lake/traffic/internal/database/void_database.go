package database

import (
	"errors"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/queue_set"
)

var (
	_ IDatabase = &sVoidDatabase{}
)

type sVoidDatabase struct {
	fSettings ISettings
	fQueueSet queue_set.IQueueSet
}

func NewVoidDatabase(pSett ISettings) IDatabase {
	return &sVoidDatabase{
		fSettings: pSett,
		fQueueSet: queue_set.NewQueueSet(
			queue_set.NewSettings(&queue_set.SSettings{
				FCapacity: pSett.GetMessagesCapacity(),
			}),
		),
	}
}

func (p *sVoidDatabase) Settings() ISettings {
	return p.fSettings
}

func (p *sVoidDatabase) Hashes() ([][]byte, error) {
	return p.fQueueSet.GetQueueKeys(), nil
}

func (p *sVoidDatabase) Push(pMsg net_message.IMessage) error {
	if gotMsg := net_message.LoadMessage(p.fSettings, pMsg.ToBytes()); gotMsg == nil {
		return errors.New("got message with diff settings")
	}

	if ok := p.fQueueSet.Push(pMsg.GetHash(), pMsg.ToBytes()); !ok {
		return GErrMessageIsExist
	}

	return nil
}

func (p *sVoidDatabase) Load(pHash []byte) (net_message.IMessage, error) {
	if len(pHash) != hashing.CSHA256Size {
		return nil, errors.New("key size invalid")
	}

	msgBytes, ok := p.fQueueSet.Load(pHash)
	if !ok {
		return nil, GErrMessageIsNotExist
	}

	msg := net_message.LoadMessage(p.fSettings, msgBytes)
	if msg == nil {
		panic("message is nil")
	}

	return msg, nil
}

func (p *sVoidDatabase) Close() error {
	return nil
}
