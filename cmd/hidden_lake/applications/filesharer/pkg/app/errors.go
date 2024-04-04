package app

import "errors"

const (
	errPrefix = "cmd/hidden_lake/applications/filesharer/pkg/app = "
)

var (
	ErrRunning = errors.New(errPrefix + "app running")
	ErrService = errors.New(errPrefix + "service")
	ErrInitSTG = errors.New(errPrefix + "init storage")
	ErrClose   = errors.New(errPrefix + "close")
)
