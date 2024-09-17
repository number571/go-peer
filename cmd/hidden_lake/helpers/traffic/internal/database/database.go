package database

import (
	"sync"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/utils"
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
	kvDB, err := database.NewKVDatabase(
		database.NewSettings(&database.SSettings{
			FPath: pSett.GetPath(),
		}),
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrCreateDB, err)
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

func (p *sDatabase) Pointer() uint64 {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.getPointer()
}

func (p *sDatabase) Hash(i uint64) ([]byte, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	messagesCapacity := p.Settings().GetMessagesCapacity()
	if i >= messagesCapacity {
		return nil, ErrGtMessagesCapacity
	}

	hash, err := p.fDB.Get(getKeyHash(i))
	if err != nil {
		return nil, utils.MergeErrors(ErrMessageIsNotExist, err)
	}

	return hash, nil
}

func (p *sDatabase) Push(pMsg net_message.IMessage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if _, err := net_message.LoadMessage(p.fSettings, pMsg.ToBytes()); err != nil {
		return utils.MergeErrors(ErrLoadMessage, err)
	}

	msgHash := pMsg.GetHash()
	if _, err := p.fDB.Get(getKeyMessage(msgHash)); err == nil {
		return utils.MergeErrors(ErrMessageIsExist, err)
	}

	keyHash := getKeyHash(p.getPointer())

	// delete old message by pointer
	if hash, err := p.fDB.Get(keyHash); err == nil {
		keyMsg := getKeyMessage(hash)
		if err := p.fDB.Del(keyMsg); err != nil {
			return utils.MergeErrors(ErrDeleteOldKey, err)
		}
	}

	// rewrite hash's field
	if err := p.fDB.Set(keyHash, msgHash); err != nil {
		return utils.MergeErrors(ErrRewriteKeyHash, err)
	}

	// write message
	keyMsg := getKeyMessage(msgHash)
	if err := p.fDB.Set(keyMsg, pMsg.ToBytes()); err != nil {
		return utils.MergeErrors(ErrWriteMessage, err)
	}

	// update pointer
	if err := p.incPointer(); err != nil {
		return utils.MergeErrors(ErrIncrementPointer, err)
	}

	return nil
}

func (p *sDatabase) Load(pHash []byte) (net_message.IMessage, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if len(pHash) != hashing.CSHA256Size {
		return nil, ErrInvalidKeySize
	}

	data, err := p.fDB.Get(getKeyMessage(pHash))
	if err != nil {
		return nil, ErrMessageIsNotExist
	}

	msg, err := net_message.LoadMessage(p.Settings(), data)
	if err != nil {
		return nil, utils.MergeErrors(ErrLoadMessage, err)
	}

	return msg, nil
}

func (p *sDatabase) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if err := p.fDB.Close(); err != nil {
		return utils.MergeErrors(ErrCloseDB, err)
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
		return utils.MergeErrors(ErrSetPointer, err)
	}
	return nil
}
