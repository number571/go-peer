package database

import "github.com/number571/go-peer/storage"

type IKeyValueDB interface {
	storage.IKeyValueStorage

	// Turn on hashing keys, encryption values
	WithHashing(bool) IKeyValueDB
	WithEncryption([]byte) IKeyValueDB

	// some databases are can be not implements this methods
	// than they are return default values Iter=nil, Close=nil
	Iter([]byte) iIterator
	Close() error
}

type iIterator interface {
	Key() []byte
	Value() []byte

	Next() bool
	Close()
}
