package app

import "errors"

var (
	ErrRunning = errors.New("app running")
	ErrService = errors.New("service")
	ErrClose   = errors.New("close")
)
