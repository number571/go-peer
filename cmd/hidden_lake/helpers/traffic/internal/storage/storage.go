package storage

import (
	"errors"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/utils"
)

var (
	_ IMessageStorage = &sMessageStorage{}
)

type sMessageStorage struct {
	fSettings net_message.ISettings
	fDatabase database.IKVDatabase
	fLRUCache cache.ILRUCache
}

func NewMessageStorage(
	pSettings net_message.ISettings,
	pDatabase database.IKVDatabase,
	pLRUCache cache.ILRUCache,
) IMessageStorage {
	return &sMessageStorage{
		fSettings: pSettings,
		fDatabase: pDatabase,
		fLRUCache: pLRUCache,
	}
}

func (p *sMessageStorage) GetSettings() net_message.ISettings {
	return p.fSettings
}

func (p *sMessageStorage) GetKVDatabase() database.IKVDatabase {
	return p.fDatabase
}

func (p *sMessageStorage) GetLRUCache() cache.ILRUCache {
	return p.fLRUCache
}

func (p *sMessageStorage) Pointer() uint64 {
	return p.fLRUCache.GetIndex()
}

func (p *sMessageStorage) Hash(i uint64) ([]byte, error) {
	key, ok := p.fLRUCache.GetKey(i)
	if !ok {
		return nil, ErrMessageIsNotExist
	}
	return key, nil
}

func (p *sMessageStorage) Push(pMsg net_message.IMessage) error {
	if _, err := net_message.LoadMessage(p.fSettings, pMsg.ToBytes()); err != nil {
		return utils.MergeErrors(ErrLoadMessage, err)
	}
	hash := pMsg.GetHash()
	_, err := p.fDatabase.Get(hash)
	if err == nil {
		return ErrHashAlreadyExist
	}
	if !errors.Is(err, database.ErrNotFound) {
		return utils.MergeErrors(ErrGetHashFromDB, err)
	}
	if err := p.fDatabase.Set(hash, []byte{}); err != nil {
		return utils.MergeErrors(ErrSetHashIntoDB, err)
	}
	if ok := p.fLRUCache.Set(hash, pMsg.ToBytes()); !ok {
		return ErrMessageIsExist
	}
	return nil
}

func (p *sMessageStorage) Load(pHash []byte) (net_message.IMessage, error) {
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
