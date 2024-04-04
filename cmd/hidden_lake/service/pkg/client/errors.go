package client

import "errors"

const (
	errPrefix = "cmd/hidden_lake/service/pkg/client = "
)

var (
	ErrBadRequest       = errors.New(errPrefix + "bad request")
	ErrDecodeResponse   = errors.New(errPrefix + "decode response")
	ErrInvalidPublicKey = errors.New(errPrefix + "invalid public key")
	ErrInvalidTitle     = errors.New(errPrefix + "invalid title")
)
