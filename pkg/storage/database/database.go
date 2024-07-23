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
	return &sKVDatabase{db}, nil
}

func (p *sKVDatabase) Set(pKey []byte, pValue []byte) error {
	err := p.fDB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(cBucket))
		if err != nil {
			return err
		}
		return bucket.Put(pKey, pValue)
	})
	if err != nil {
		return utils.MergeErrors(ErrSetValue, err)
	}
	return nil
}

func (p *sKVDatabase) Get(pKey []byte) ([]byte, error) {
	var encValue []byte
	err := p.fDB.View(func(tx *bbolt.Tx) error {
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
	if err != nil {
		return nil, utils.MergeErrors(ErrGetValue, err)
	}
	return encValue, nil
}

func (p *sKVDatabase) Del(pKey []byte) error {
	err := p.fDB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(cBucket))
		if b == nil {
			return nil
		}
		return b.Delete(pKey)
	})
	if err != nil {
		return utils.MergeErrors(ErrDelValue, err)
	}
	return nil
}

func (p *sKVDatabase) Close() error {
	if err := p.fDB.Close(); err != nil {
		return utils.MergeErrors(ErrCloseDB, err)
	}
	return nil
}
