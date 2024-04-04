package app

import "errors"

var (
	ErrRunning           = errors.New("app running")
	ErrService           = errors.New("service")
	ErrInitDB            = errors.New("init database")
	ErrClose             = errors.New("close")
	ErrInvalidPrivateKey = errors.New("invalid private key")
	ErrReadPrivateKey    = errors.New("read private key")
	ErrWritePrivateKey   = errors.New("write private key")
	ErrSizePrivateKey    = errors.New("size private key")
	ErrGetPrivateKey     = errors.New("get private key")
	ErrInitConfig        = errors.New("init config")
	ErrSetParallelNull   = errors.New("set parallel = 0")
	ErrGetParallel       = errors.New("get parallel")
)
