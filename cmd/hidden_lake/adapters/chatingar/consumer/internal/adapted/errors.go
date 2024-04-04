package adapted

import "errors"

const (
	errPrefix = "cmd/hidden_lake/adapters/chatingar/consumer/internal/adapted = "
)

var (
	ErrBuildRequest      = errors.New(errPrefix + "build request")
	ErrBadRequest        = errors.New(errPrefix + "bad request")
	ErrBadStatusCode     = errors.New(errPrefix + "bad status code")
	ErrDecodeCount       = errors.New(errPrefix + "decode count")
	ErrCountLtNull       = errors.New(errPrefix + "count < 0")
	ErrLimitPage         = errors.New(errPrefix + "limit page")
	ErrDecodeMessages    = errors.New(errPrefix + "ldecode messages")
	ErrLoadCountComments = errors.New(errPrefix + "load count comments")
)
