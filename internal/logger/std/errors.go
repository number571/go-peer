package std

const (
	errPrefix = "internal/logger/std = "
)

type SStdError struct {
	str string
}

func (err *SStdError) Error() string { return errPrefix + err.str }

var (
	ErrUnknownLogType = &SStdError{"unknown log type"}
)
