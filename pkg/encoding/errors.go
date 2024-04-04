package encoding

import "errors"

const (
	errPrefix = "pkg/encoding = "
)

var (
	ErrDeserialize = errors.New(errPrefix + "deserialize bytes")
)
