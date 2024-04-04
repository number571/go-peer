package app

import "errors"

const (
	errPrefix = "cmd/hidden_lake/service/pkg/app = "
)

var (
	ErrRunning           = errors.New(errPrefix + "app running")
	ErrService           = errors.New(errPrefix + "service")
	ErrInitDB            = errors.New(errPrefix + "init database")
	ErrClose             = errors.New(errPrefix + "close")
	ErrInvalidPrivateKey = errors.New(errPrefix + "invalid private key")
	ErrReadPrivateKey    = errors.New(errPrefix + "read private key")
	ErrWritePrivateKey   = errors.New(errPrefix + "write private key")
	ErrSizePrivateKey    = errors.New(errPrefix + "size private key")
	ErrGetPrivateKey     = errors.New(errPrefix + "get private key")
	ErrInitConfig        = errors.New(errPrefix + "init config")
	ErrSetParallelNull   = errors.New(errPrefix + "set parallel = 0")
	ErrGetParallel       = errors.New(errPrefix + "get parallel")
)
