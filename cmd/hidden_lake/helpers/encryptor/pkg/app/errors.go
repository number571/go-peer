package app

import "errors"

var (
	ErrRunning           = errors.New("app running")
	ErrService           = errors.New("service")
	ErrClose             = errors.New("close")
	ErrSizePrivateKey    = errors.New("size private key")
	ErrGetPrivateKey     = errors.New("get private key")
	ErrInitConfig        = errors.New("init config")
	ErrSetParallelNull   = errors.New("set parallel = 0")
	ErrGetParallelValue  = errors.New("get parallel value")
	ErrWritePrivateKey   = errors.New("write private key")
	ErrReadPrivateKey    = errors.New("read private key")
	ErrInvalidPrivateKey = errors.New("invalid private key")
)
