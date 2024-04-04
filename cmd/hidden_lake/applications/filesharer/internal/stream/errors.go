package stream

import "errors"

const (
	errPrefix = "cmd/hidden_lake/applications/filesharer/internal/stream = "
)

var (
	ErrWriteFileChunk   = errors.New(errPrefix + "write file chunk")
	ErrLoadFileChunk    = errors.New(errPrefix + "load file chunk")
	ErrInvalidHash      = errors.New(errPrefix + "invalid hash")
	ErrRetryFailed      = errors.New(errPrefix + "retry failed")
	ErrInvalidWhence    = errors.New(errPrefix + "invalid whence")
	ErrNegativePosition = errors.New(errPrefix + "negative position")
)
