package storage

import "errors"

var (
	ErrReadFile         = errors.New("read file")
	ErrWriteFile        = errors.New("write file")
	ErrSaltSize         = errors.New("size of storage < salt size")
	ErrInitStorage      = errors.New("set init storage")
	ErrDecryptStorage   = errors.New("decrypt storage")
	ErrStorageUndefined = errors.New("storage undefined")
	ErrKeyIsNotExist    = errors.New("key is not exist")
	ErrUnmarshalMap     = errors.New("unmarshal decrypted map")
)
