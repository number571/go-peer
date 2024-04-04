package app

import "errors"

const (
	errPrefix = "cmd/hidden_lake/composite/pkg/app = "
)

var (
	ErrRunning        = errors.New(errPrefix + "app running")
	ErrService        = errors.New(errPrefix + "service")
	ErrClose          = errors.New(errPrefix + "close")
	ErrUnknownService = errors.New(errPrefix + "unknown service")
	ErrGetRunners     = errors.New(errPrefix + "get runners")
	ErrInitConfig     = errors.New(errPrefix + "init config")
)
