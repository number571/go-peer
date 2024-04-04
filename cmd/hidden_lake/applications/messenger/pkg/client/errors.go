package client

import "errors"

const (
	errPrefix = "cmd/hidden_lake/applications/messenger/pkg/client = "
)

var (
	ErrBadRequest     = errors.New(errPrefix + "bad request")
	ErrDecodeResponse = errors.New(errPrefix + "decode response")
)
