package adapted

import "errors"

const (
	errPrefix = "cmd/hidden_lake/adapters/common/producer/internal/adapted = "
)

var (
	ErrInvalidResponse = errors.New(errPrefix + "invalid response")
	ErrReadResponse    = errors.New(errPrefix + "read response")
	ErrBadStatusCode   = errors.New(errPrefix + "bad status code")
	ErrBadRequest      = errors.New(errPrefix + "bad request")
	ErrBuildRequest    = errors.New(errPrefix + "build request")
)
