package client

import "errors"

const (
	errPrefix = "cmd/hidden_lake/helpers/loader/pkg/client = "
)

var (
	ErrBadRequest     = errors.New(errPrefix + "bad request")
	ErrDecodeResponse = errors.New(errPrefix + "decode response")
	ErrInvalidTitle   = errors.New(errPrefix + "invalid title")
)
