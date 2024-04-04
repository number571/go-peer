package std

import "errors"

const (
	errPrefix = "internal/logger/std = "
)

var (
	ErrUnknownLogType = errors.New(errPrefix + "unknown log type")
)
