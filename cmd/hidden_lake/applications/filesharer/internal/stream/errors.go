package stream

const (
	errPrefix = "cmd/hidden_lake/applications/filesharer/internal/stream = "
)

type SStreamError struct {
	str string
}

func (err *SStreamError) Error() string {
	return errPrefix + err.str
}

var (
	ErrWriteFileChunk   = &SStreamError{"write file chunk"}
	ErrLoadFileChunk    = &SStreamError{"load file chunk"}
	ErrInvalidHash      = &SStreamError{"invalid hash"}
	ErrRetryFailed      = &SStreamError{"retry failed"}
	ErrInvalidWhence    = &SStreamError{"invalid whence"}
	ErrNegativePosition = &SStreamError{"negative position"}
	ErrGetMessageLimit  = &SStreamError{"get message limit"}
)
