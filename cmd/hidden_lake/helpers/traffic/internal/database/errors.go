package database

import "github.com/syndtr/goleveldb/leveldb/errors"

var (
	ErrMessageIsExist    = errors.New("message is exist")
	ErrMessageIsNotExist = errors.New("message is not exist")
)
