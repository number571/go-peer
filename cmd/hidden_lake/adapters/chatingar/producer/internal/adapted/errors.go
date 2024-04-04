package adapted

import "errors"

const (
	errPrefix = "cmd/hidden_lake/adapters/chatingar/producer/internal/adapted = "
)

var (
	ErrBadStatusCode = errors.New(errPrefix + "bad status code")
	ErrBadRequest    = errors.New(errPrefix + "bad request")
	ErrBuildRequest  = errors.New(errPrefix + "build request")
)
