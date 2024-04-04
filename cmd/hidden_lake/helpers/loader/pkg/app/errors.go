package app

import "errors"

const (
	errPrefix = "cmd/hidden_lake/helpers/loader/pkg/app = "
)

var (
	ErrRunning = errors.New(errPrefix + "app running")
	ErrService = errors.New(errPrefix + "service")
	ErrClose   = errors.New(errPrefix + "close")
)
