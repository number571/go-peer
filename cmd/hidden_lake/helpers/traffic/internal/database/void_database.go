package database

import "github.com/number571/go-peer/pkg/storage/database"

var _ database.IKVDatabase = &sVoidKVDatabase{}

type sVoidKVDatabase struct{}

func NewVoidKVDatabase() database.IKVDatabase {
	return &sVoidKVDatabase{}
}

func (p *sVoidKVDatabase) Set([]byte, []byte) error   { return nil }
func (p *sVoidKVDatabase) Get([]byte) ([]byte, error) { return nil, database.ErrNotFound }
func (p *sVoidKVDatabase) Del([]byte) error           { return nil }
func (p *sVoidKVDatabase) Close() error               { return nil }
