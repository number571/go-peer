package app

import "errors"

var (
	ErrRunning        = errors.New("app running")
	ErrService        = errors.New("service")
	ErrClose          = errors.New("close")
	ErrUnknownService = errors.New("unknown service")
	ErrGetRunners     = errors.New("get runners")
	ErrInitConfig     = errors.New("init config")
)
