package app

import "errors"

const (
	errPrefix = "cmd/hidden_lake/helpers/traffic/pkg/app = "
)

var (
	ErrRunning = errors.New(errPrefix + "app running")
	ErrService = errors.New(errPrefix + "service")
	ErrInitDB  = errors.New(errPrefix + "init database")
	ErrClose   = errors.New(errPrefix + "close")
)
