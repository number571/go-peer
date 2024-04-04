package app

import "errors"

const (
	errPrefix = "cmd/hidden_lake/helpers/encryptor/pkg/app = "
)

var (
	ErrRunning           = errors.New(errPrefix + "app running")
	ErrService           = errors.New(errPrefix + "service")
	ErrClose             = errors.New(errPrefix + "close")
	ErrSizePrivateKey    = errors.New(errPrefix + "size private key")
	ErrGetPrivateKey     = errors.New(errPrefix + "get private key")
	ErrInitConfig        = errors.New(errPrefix + "init config")
	ErrSetParallelNull   = errors.New(errPrefix + "set parallel = 0")
	ErrGetParallelValue  = errors.New(errPrefix + "get parallel value")
	ErrWritePrivateKey   = errors.New(errPrefix + "write private key")
	ErrReadPrivateKey    = errors.New(errPrefix + "read private key")
	ErrInvalidPrivateKey = errors.New(errPrefix + "invalid private key")
)
