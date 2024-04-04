package client

import "errors"

const (
	errPrefix = "cmd/hidden_lake/helpers/traffic/pkg/client = "
)

var (
	ErrBadRequest       = errors.New(errPrefix + "bad request")
	ErrDecodeResponse   = errors.New(errPrefix + "decode response")
	ErrInvalidResponse  = errors.New(errPrefix + "invalid response")
	ErrInvalidHexFormat = errors.New(errPrefix + "invalid hex format")
	ErrInvalidTitle     = errors.New(errPrefix + "invalid title")
	ErrDecodeMessage    = errors.New(errPrefix + "decode message")
)
