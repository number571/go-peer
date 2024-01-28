package app

import "errors"

var (
	ErrRunning = errors.New("app running")
	ErrService = errors.New("service")
	ErrInitDB  = errors.New("init database")
	ErrClose   = errors.New("close")
)
