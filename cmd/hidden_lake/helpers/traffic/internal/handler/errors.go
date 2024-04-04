package handler

import "errors"

const (
	errPrefix = "cmd/hidden_lake/helpers/traffic/internal/handler = "
)

var (
	ErrLoadMessage   = errors.New(errPrefix + "load message")
	ErrDatabaseNull  = errors.New(errPrefix + "database null")
	ErrPushMessageDB = errors.New(errPrefix + "push message db")
)
