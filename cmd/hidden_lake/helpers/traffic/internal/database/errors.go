package database

import "github.com/syndtr/goleveldb/leveldb/errors"

var (
	ErrMessageIsExist    = errors.New("message is exist")
	ErrMessageIsNotExist = errors.New("message is not exist")
	ErrInvalidKeySize    = errors.New("invalid key size")
	ErrLoadMessage       = errors.New("load message")
	ErrCloseDB           = errors.New("close db")
	ErrSetPointer        = errors.New("set pointer")
	ErrIncrementPointer  = errors.New("increment pointer")
	ErrWriteMessage      = errors.New("write message")
	ErrRewriteKeyHash    = errors.New("rewrite key hash")
	ErrDeleteOldKey      = errors.New("delete old key")
	ErrCreateDB          = errors.New("create db")
)
