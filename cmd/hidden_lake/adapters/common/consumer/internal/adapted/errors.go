package adapted

import "errors"

const (
	errPrefix = "cmd/hidden_lake/adapters/common/consumer/internal/adapted = "
)

var (
	ErrLoadCountService = errors.New(errPrefix + "load count from service")
	ErrLoadCountDB      = errors.New(errPrefix + "load count from db")
	ErrGetCount         = errors.New(errPrefix + "get count")
	ErrSetNewCount      = errors.New(errPrefix + "set new count")
	ErrInitCountKey     = errors.New(errPrefix + "init count key")
	ErrParseCount       = errors.New(errPrefix + "parse count")
	ErrInvalidResponse  = errors.New(errPrefix + "invalid response")
	ErrReadResponse     = errors.New(errPrefix + "read response")
	ErrBadRequest       = errors.New(errPrefix + "bad request")
	ErrBuildRequest     = errors.New(errPrefix + "build request")
	ErrDecodeMessage    = errors.New(errPrefix + "decode message")
	ErrIncrementCount   = errors.New(errPrefix + "increment count")
	ErrLoadMessage      = errors.New(errPrefix + "load message")
)
