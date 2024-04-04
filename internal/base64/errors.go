package base64

import "errors"

const (
	errPrefix = "internal/base64 = "
)

var (
	ErrBase64Size = errors.New(errPrefix + "base64 size")
)
