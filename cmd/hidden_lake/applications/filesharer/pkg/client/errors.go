package client

import "errors"

const (
	errPrefix = "cmd/hidden_lake/applications/filesharer/pkg/client = "
)

var (
	ErrBadRequest      = errors.New(errPrefix + "bad request")
	ErrDecodeResponse  = errors.New(errPrefix + "decode response")
	ErrInvalidResponse = errors.New(errPrefix + "invalid response")
)
