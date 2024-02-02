package app

import "errors"

var (
	ErrRunning = errors.New("app running")
	ErrService = errors.New("service")
	ErrInitSTG = errors.New("init storage")
	ErrClose   = errors.New("close")
)
