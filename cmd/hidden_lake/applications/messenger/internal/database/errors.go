package database

import "errors"

var (
	ErrLoadMessage    = errors.New("load message")
	ErrGetMessage     = errors.New("get message")
	ErrSetMessage     = errors.New("set message")
	ErrSetSizeMessage = errors.New("set size message")
	ErrCloseDB        = errors.New("close db")
	ErrEndGtSize      = errors.New("end > size")
	ErrStartGtEnd     = errors.New("start > end")
	ErrCreateDB       = errors.New("create db")
)
