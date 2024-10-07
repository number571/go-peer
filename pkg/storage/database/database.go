package database

import (
	"github.com/number571/go-peer/pkg/utils"
	"go.etcd.io/bbolt"
)

var (
	_ IKVDatabase = &sKVDatabase{}
)

const (
	cBucket = "_BUCKET_"
)

type sKVDatabase struct {
	fDB *bbolt.DB
}

func NewKVDatabase(pPath string) (IKVDatabase, error) {
	db, err := bbolt.Open(pPath, 0600, &bbolt.Options{})
	if err != nil {
		return nil, utils.MergeErrors(ErrOpenDB, err)
	}

	return &sKVDatabase{fDB: db}, nil
}

func (p *sKVDatabase) Set(pKey []byte, pValue []byte) error {
	if err := setDB(p.fDB, pKey, pValue); err != nil {
		return utils.MergeErrors(ErrSetValueDB, err)
	}
	return nil
}

func setDB(pDB *bbolt.DB, pKey []byte, pValue []byte) error {
	return pDB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(cBucket))
		if err != nil {
			return err
		}
		return bucket.Put(pKey, pValue)
	})
}

func (p *sKVDatabase) Get(pKey []byte) ([]byte, error) {
	value, err := getDB(p.fDB, pKey)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetValueDB, err)
	}
	return value, nil
}

func getDB(pDB *bbolt.DB, pKey []byte) ([]byte, error) {
	var encValue []byte
	err := pDB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(cBucket))
		if b == nil {
			return ErrNotFound
		}
		val := b.Get(pKey)
		if val == nil {
			return ErrNotFound
		}
		encValue = make([]byte, len(val))
		copy(encValue, val)
		return nil
	})
	return encValue, err
}

func (p *sKVDatabase) Del(pKey []byte) error {
	if err := delDB(p.fDB, pKey); err != nil {
		return utils.MergeErrors(ErrDelValueDB, err)
	}
	return nil
}

func delDB(pDB *bbolt.DB, pKey []byte) error {
	return pDB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(cBucket))
		if b == nil {
			return nil
		}
		return b.Delete(pKey)
	})
}

func (p *sKVDatabase) Close() error {
	if err := p.fDB.Close(); err != nil {
		return utils.MergeErrors(ErrCloseDB, err)
	}
	return nil
}
