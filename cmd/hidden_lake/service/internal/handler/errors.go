package handler

import "errors"

const (
	errPrefix = "cmd/hidden_lake/service/internal/handler = "
)

var (
	ErrBadRequest          = errors.New(errPrefix + "bad request")
	ErrBuildRequest        = errors.New(errPrefix + "build request")
	ErrUndefinedService    = errors.New(errPrefix + "undefined service")
	ErrLoadRequest         = errors.New(errPrefix + "load request")
	ErrUpdateFriends       = errors.New(errPrefix + "update friends")
	ErrInvalidResponseMode = errors.New(errPrefix + "invalid response mode")
)
