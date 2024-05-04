package app

const (
	errPrefix = "cmd/hidden_lake/service/pkg/app = "
)

type SAppError struct {
	str string
}

func (err *SAppError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRunning           = &SAppError{"app running"}
	ErrService           = &SAppError{"service"}
	ErrInitDB            = &SAppError{"init database"}
	ErrClose             = &SAppError{"close"}
	ErrInvalidPrivateKey = &SAppError{"invalid private key"}
	ErrReadPrivateKey    = &SAppError{"read private key"}
	ErrWritePrivateKey   = &SAppError{"write private key"}
	ErrSizePrivateKey    = &SAppError{"size private key"}
	ErrGetPrivateKey     = &SAppError{"get private key"}
	ErrInitConfig        = &SAppError{"init config"}
	ErrSetParallelNull   = &SAppError{"set parallel = 0"}
	ErrGetParallel       = &SAppError{"get parallel"}
	ErrCreateAnonNode    = &SAppError{"create anon node"}
	ErrOpenKVDatabase    = &SAppError{"open kv database"}
)
