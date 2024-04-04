package database

import "errors"

const (
	errPrefix = "pkg/database = "
)

var (
	ErrOpenDB               = errors.New(errPrefix + "open database")
	ErrReadSalt             = errors.New(errPrefix + "read salt value")
	ErrReadSaltHash         = errors.New(errPrefix + "read salt hash")
	ErrPushSalt             = errors.New(errPrefix + "push salt value")
	ErrPushSaltHash         = errors.New(errPrefix + "push salt hash")
	ErrInvalidSaltHash      = errors.New(errPrefix + "invalid salt hash")
	ErrSetValueDB           = errors.New(errPrefix + "set value to database")
	ErrGetValueDB           = errors.New(errPrefix + "get value from database")
	ErrDelValueDB           = errors.New(errPrefix + "del value from database")
	ErrCloseDB              = errors.New(errPrefix + "close database")
	ErrRecoverDB            = errors.New(errPrefix + "recover database")
	ErrInvalidEncryptedSize = errors.New(errPrefix + "invalid encrypted size")
	ErrInvalidDataHash      = errors.New(errPrefix + "invalid data hash")
)
