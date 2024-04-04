package handler

import "errors"

var (
	ErrBadRequest          = errors.New("bad request")
	ErrBuildRequest        = errors.New("build request")
	ErrUndefinedService    = errors.New("undefined service")
	ErrLoadRequest         = errors.New("load request")
	ErrUpdateFriends       = errors.New("update friends")
	ErrInvalidResponseMode = errors.New("invalid response mode")
)
