package database

import "errors"

var (
	ErrOpenDB               = errors.New("open database")
	ErrReadSalt             = errors.New("read salt value")
	ErrReadSaltHash         = errors.New("read salt hash")
	ErrPushSalt             = errors.New("push salt value")
	ErrPushSaltHash         = errors.New("push salt hash")
	ErrInvalidSaltHash      = errors.New("invalid salt hash")
	ErrSetValueDB           = errors.New("set value to database")
	ErrGetValueDB           = errors.New("get value from database")
	ErrDelValueDB           = errors.New("del value from database")
	ErrCloseDB              = errors.New("close database")
	ErrRecoverDB            = errors.New("recover database")
	ErrInvalidEncryptedSize = errors.New("invalid encrypted size")
	ErrInvalidDataHash      = errors.New("invalid data hash")
)
