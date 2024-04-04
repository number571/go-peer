package stream

import "errors"

var (
	ErrWriteFileChunk   = errors.New("write file chunk")
	ErrLoadFileChunk    = errors.New("load file chunk")
	ErrInvalidHash      = errors.New("invalid hash")
	ErrRetryFailed      = errors.New("retry failed")
	ErrInvalidWhence    = errors.New("invalid whence")
	ErrNegativePosition = errors.New("negative position")
)
