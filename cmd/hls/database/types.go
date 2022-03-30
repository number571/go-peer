package database

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type IKeyValueDB interface {
	Push([]byte) error
	Exist([]byte) bool

	Close() error
	Clean() error

	dbPointer() *leveldb.DB
}
