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
	ErrRunning          = &SAppError{"app running"}
	ErrService          = &SAppError{"service"}
	ErrInitDB           = &SAppError{"init database"}
	ErrClose            = &SAppError{"close"}
	ErrSizePrivateKey   = &SAppError{"size private key"}
	ErrGetPrivateKey    = &SAppError{"get private key"}
	ErrInitConfig       = &SAppError{"init config"}
	ErrSetParallelNull  = &SAppError{"set parallel = 0"}
	ErrGetParallel      = &SAppError{"get parallel"}
	ErrCreateAnonNode   = &SAppError{"create anon node"}
	ErrOpenKVDatabase   = &SAppError{"open kv database"}
	ErrReadKVDatabase   = &SAppError{"read kv database"}
	ErrMessageSizeLimit = &SAppError{"message size limit"}
	ErrInvalidPsdPubKey = &SAppError{"invalid psd public key"}
	ErrGetPsdPubKey     = &SAppError{"get psd pub key"}
	ErrSetPsdPubKey     = &SAppError{"set psd pub key"}
)
