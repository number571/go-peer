package database

import "github.com/syndtr/goleveldb/leveldb/errors"

var (
	GErrMessageIsExist    = errors.New("message is exist")
	GErrMessageIsNotExist = errors.New("message is not exist")
)
