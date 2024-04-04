package database

import "github.com/syndtr/goleveldb/leveldb/errors"

const (
	errPrefix = "cmd/hidden_lake/helpers/traffic/internal/database = "
)

var (
	ErrMessageIsExist    = errors.New(errPrefix + "message is exist")
	ErrMessageIsNotExist = errors.New(errPrefix + "message is not exist")
	ErrInvalidKeySize    = errors.New(errPrefix + "invalid key size")
	ErrLoadMessage       = errors.New(errPrefix + "load message")
	ErrCloseDB           = errors.New(errPrefix + "close db")
	ErrSetPointer        = errors.New(errPrefix + "set pointer")
	ErrIncrementPointer  = errors.New(errPrefix + "increment pointer")
	ErrWriteMessage      = errors.New(errPrefix + "write message")
	ErrRewriteKeyHash    = errors.New(errPrefix + "rewrite key hash")
	ErrDeleteOldKey      = errors.New(errPrefix + "delete old key")
	ErrCreateDB          = errors.New(errPrefix + "create db")
)
