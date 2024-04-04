package handler

import "errors"

var (
	ErrLoadMessage   = errors.New("load message")
	ErrDatabaseNull  = errors.New("database null")
	ErrPushMessageDB = errors.New("push message db")
)
