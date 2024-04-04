package database

import "errors"

const (
	errPrefix = "cmd/hidden_lake/applications/messenger/internal/database = "
)

var (
	ErrLoadMessage    = errors.New(errPrefix + "load message")
	ErrGetMessage     = errors.New(errPrefix + "get message")
	ErrSetMessage     = errors.New(errPrefix + "set message")
	ErrSetSizeMessage = errors.New(errPrefix + "set size message")
	ErrCloseDB        = errors.New(errPrefix + "close db")
	ErrEndGtSize      = errors.New(errPrefix + "end > size")
	ErrStartGtEnd     = errors.New(errPrefix + "start > end")
	ErrCreateDB       = errors.New(errPrefix + "create db")
)
