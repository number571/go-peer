package database

import (
	"errors"
	"fmt"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage"
	"github.com/number571/go-peer/pkg/storage/database"
)

var (
	_ IDatabase = &sDatabase{}
)

type sDatabase struct {
	fMutex   sync.Mutex
	fPointer uint64

	fSettings ISettings
	fDB       database.IKVDatabase
}

func NewDatabase(pSett ISettings) (IDatabase, error) {
	kvDB, err := database.NewKeyValueDB(
		storage.NewSettings(&storage.SSettings{
			FPath: pSett.GetPath(),
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("new key/value database: %w", err)
	}

	db := &sDatabase{
		fSettings: pSett,
		fDB:       kvDB,
	}
	db.fPointer = db.getPointer()
	return db, nil
}

func (p *sDatabase) Settings() ISettings {
	return p.fSettings
}

func (p *sDatabase) Hashes() ([][]byte, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	msgsLimit := p.Settings().GetMessagesCapacity()
	res := make([][]byte, 0, msgsLimit)
	for i := uint64(0); i < msgsLimit; i++ {
		hash, err := p.fDB.Get(getKeyHash(i))
		if err != nil {
			break
		}
		if len(hash) != hashing.CSHA256Size {
			panic("incorrect hash size")
		}
		res = append(res, hash)
	}

	return res, nil
}

func (p *sDatabase) Push(pMsg net_message.IMessage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if _, err := net_message.LoadMessage(p.fSettings, pMsg.ToBytes()); err != nil {
		return errors.New("got message with diff settings")
	}

	msgHash := pMsg.GetHash()
	if _, err := p.fDB.Get(getKeyMessage(msgHash)); err == nil {
		return GErrMessageIsExist
	}

	keyHash := getKeyHash(p.getPointer())

	// delete old message by pointer
	if hash, err := p.fDB.Get(keyHash); err == nil {
		keyMsg := getKeyMessage(hash)
		if err := p.fDB.Del(keyMsg); err != nil {
			return fmt.Errorf("delete old key: %w", err)
		}
	}

	// rewrite hash's field
	if err := p.fDB.Set(keyHash, msgHash); err != nil {
		return fmt.Errorf("rewrite key hash: %w", err)
	}

	// write message
	keyMsg := getKeyMessage(msgHash)
	if err := p.fDB.Set(keyMsg, pMsg.ToBytes()); err != nil {
		return fmt.Errorf("write message: %w", err)
	}

	// update pointer
	if err := p.incPointer(); err != nil {
		return fmt.Errorf("increment pointer: %w", err)
	}

	return nil
}

func (p *sDatabase) Load(pHash []byte) (net_message.IMessage, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if len(pHash) != hashing.CSHA256Size {
		return nil, errors.New("key size invalid")
	}

	data, err := p.fDB.Get(getKeyMessage(pHash))
	if err != nil {
		return nil, GErrMessageIsNotExist
	}

	msg, err := net_message.LoadMessage(p.Settings(), data)
	if err != nil {
		return nil, fmt.Errorf("load message: %w", err)
	}

	return msg, nil
}

func (p *sDatabase) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if err := p.fDB.Close(); err != nil {
		return fmt.Errorf("close KV database: %w", err)
	}
	return nil
}

func (p *sDatabase) getPointer() uint64 {
	data, err := p.fDB.Get(getKeyPointer())
	if err != nil {
		return 0
	}

	res := [encoding.CSizeUint64]byte{}
	copy(res[:], data)
	return encoding.BytesToUint64(res)
}

func (p *sDatabase) incPointer() error {
	msgsLimit := p.Settings().GetMessagesCapacity()
	res := encoding.Uint64ToBytes((p.getPointer() + 1) % msgsLimit)
	if err := p.fDB.Set(getKeyPointer(), res[:]); err != nil {
		return fmt.Errorf("set pointer into KV database: %w", err)
	}
	return nil
}
