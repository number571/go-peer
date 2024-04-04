package client

import "errors"

const (
	errPrefix = "cmd/hidden_lake/helpers/encryptor/pkg/client = "
)

var (
	ErrRequest          = errors.New(errPrefix + "request")
	ErrDecodeResponse   = errors.New(errPrefix + "decode response")
	ErrInvalidHexFormat = errors.New(errPrefix + "invalid hex format")
	ErrInvalidPublicKey = errors.New(errPrefix + "invalid public key")
	ErrInvalidTitle     = errors.New(errPrefix + "invalid title")
)
