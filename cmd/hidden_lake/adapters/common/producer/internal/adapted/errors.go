package adapted

import "errors"

var (
	ErrInvalidResponse = errors.New("invalid response")
	ErrReadResponse    = errors.New("read response")
	ErrBadStatusCode   = errors.New("bad status code")
	ErrBadRequest      = errors.New("bad request")
	ErrBuildRequest    = errors.New("build request")
)
