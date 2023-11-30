package database

import (
	"errors"
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/queue_set"
)

var (
	_ IDatabase = &sInMemoryDatabase{}
)

type sInMemoryDatabase struct {
	fSettings ISettings
	fQueueSet queue_set.IQueueSet
}

func NewInMemoryDatabase(pSett ISettings) IDatabase {
	return &sInMemoryDatabase{
		fSettings: pSett,
		fQueueSet: queue_set.NewQueueSet(
			queue_set.NewSettings(&queue_set.SSettings{
				FCapacity: pSett.GetMessagesCapacity(),
			}),
		),
	}
}

func (p *sInMemoryDatabase) Settings() ISettings {
	return p.fSettings
}

func (p *sInMemoryDatabase) Hashes() ([][]byte, error) {
	return p.fQueueSet.GetQueueKeys(), nil
}

func (p *sInMemoryDatabase) Push(pMsg net_message.IMessage) error {
	if _, err := net_message.LoadMessage(p.fSettings, pMsg.ToBytes()); err != nil {
		return errors.New("got message with diff settings")
	}

	if ok := p.fQueueSet.Push(pMsg.GetHash(), pMsg.ToBytes()); !ok {
		return GErrMessageIsExist
	}

	return nil
}

func (p *sInMemoryDatabase) Load(pHash []byte) (net_message.IMessage, error) {
	if len(pHash) != hashing.CSHA256Size {
		return nil, errors.New("key size invalid")
	}

	msgBytes, ok := p.fQueueSet.Load(pHash)
	if !ok {
		return nil, GErrMessageIsNotExist
	}

	msg, err := net_message.LoadMessage(p.fSettings, msgBytes)
	if err != nil {
		return nil, fmt.Errorf("load message: %w", err)
	}

	return msg, nil
}

func (p *sInMemoryDatabase) Close() error {
	return nil
}
