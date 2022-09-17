package database

import "github.com/number571/go-peer/modules/crypto/asymmetric"

type IRelation interface {
	Friend() asymmetric.IPubKey
	IAm() asymmetric.IPubKey
}

type IKeyValueDB interface {
	Size(IRelation) (uint64, error)
	Push(IRelation, string) error
	Load(IRelation, uint64, uint64) ([]string, error)

	Close() error
}
