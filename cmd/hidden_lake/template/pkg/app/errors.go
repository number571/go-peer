package app

const (
	errPrefix = "cmd/hidden_lake/template/pkg/app = "
)

type SAppError struct {
	str string
}

func (err *SAppError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRunning = &SAppError{"app running"}
	ErrService = &SAppError{"service"}
	ErrClose   = &SAppError{"close"}
)
