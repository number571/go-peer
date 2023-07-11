package database

import (
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/errors"
)

type sVoidKeyValueDB struct {
	fSettings ISettings
}

func NewVoidKeyValueDB(pSett ISettings) IKVDatabase {
	return &sVoidKeyValueDB{
		fSettings: pSett,
	}
}

func (p *sVoidKeyValueDB) Settings() ISettings {
	return p.fSettings
}

func (p *sVoidKeyValueDB) Hashes() ([]string, error) {
	return nil, nil
}

func (p *sVoidKeyValueDB) Push(pMsg message.IMessage) error {
	return nil
}

func (p *sVoidKeyValueDB) Load(pStrHash string) (message.IMessage, error) {
	return nil, errors.NewError("message undefined")
}

func (p *sVoidKeyValueDB) Close() error {
	return nil
}
