package database

import (
	"github.com/number571/go-peer/pkg/cache/lru"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	_ IDatabase = &sInMemoryDatabase{}
)

type sInMemoryDatabase struct {
	fSettings ISettings
	fLRUCache lru.ILRUCache
}

func NewInMemoryDatabase(pSett ISettings) (IDatabase, error) {
	return &sInMemoryDatabase{
		fSettings: pSett,
		fLRUCache: lru.NewLRUCache(
			lru.NewSettings(&lru.SSettings{
				FCapacity: pSett.GetMessagesCapacity(),
			}),
		),
	}, nil
}

func (p *sInMemoryDatabase) Settings() ISettings {
	return p.fSettings
}

func (p *sInMemoryDatabase) Pointer() uint64 {
	return p.fLRUCache.GetIndex()
}

func (p *sInMemoryDatabase) Hash(i uint64) ([]byte, error) {
	key, ok := p.fLRUCache.GetKey(i)
	if !ok {
		return nil, ErrMessageIsNotExist
	}
	return key, nil
}

func (p *sInMemoryDatabase) Push(pMsg net_message.IMessage) error {
	if _, err := net_message.LoadMessage(p.fSettings, pMsg.ToBytes()); err != nil {
		return utils.MergeErrors(ErrLoadMessage, err)
	}

	if ok := p.fLRUCache.Set(pMsg.GetHash(), pMsg.ToBytes()); !ok {
		return ErrMessageIsExist
	}

	return nil
}

func (p *sInMemoryDatabase) Load(pHash []byte) (net_message.IMessage, error) {
	if len(pHash) != hashing.CSHA256Size {
		return nil, ErrInvalidKeySize
	}

	msgBytes, ok := p.fLRUCache.Get(pHash)
	if !ok {
		return nil, ErrMessageIsNotExist
	}

	msg, err := net_message.LoadMessage(p.fSettings, msgBytes)
	if err != nil {
		return nil, utils.MergeErrors(ErrLoadMessage, err)
	}

	return msg, nil
}

func (p *sInMemoryDatabase) Close() error {
	return nil
}
