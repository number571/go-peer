package api

import "errors"

const (
	errPrefix = "internal/api = "
)

var (
	ErrBadStatusCode = errors.New(errPrefix + "bad status code")
	ErrReadResponse  = errors.New(errPrefix + "read response")
	ErrLoadResponse  = errors.New(errPrefix + "load response")
	ErrBadRequest    = errors.New(errPrefix + "bad request")
	ErrBuildRequest  = errors.New(errPrefix + "build request")
	ErrCopyBytes     = errors.New(errPrefix + "copy bytes")
)
