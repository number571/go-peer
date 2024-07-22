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
	fSettings ISettings
	fDB       *bbolt.DB
}

func NewKVDatabase(pSett ISettings) (IKVDatabase, error) {
	db, err := bbolt.Open(pSett.GetPath(), 0600, &bbolt.Options{})
	if err != nil {
		return nil, utils.MergeErrors(ErrOpenDB, err)
	}
	return &sKVDatabase{
		fSettings: pSett,
		fDB:       db,
	}, nil
}

func (p *sKVDatabase) GetSettings() ISettings {
	return p.fSettings
}

func (p *sKVDatabase) Set(pKey []byte, pValue []byte) error {
	return p.fDB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(cBucket))
		if err != nil {
			return err
		}
		return bucket.Put(pKey, pValue)
	})
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
	return encValue, err
}

func (p *sKVDatabase) Del(pKey []byte) error {
	return p.fDB.Update(func(tx *bbolt.Tx) error {
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
