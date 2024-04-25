package joiner

const (
	errPrefix = "pkg/payload/joiner = "
)

type SJoinerError struct {
	str string
}

func (err *SJoinerError) Error() string {
	return errPrefix + err.str
}

var (
	ErrLoadPayload   = &SJoinerError{"load payload"}
	ErrInvalidLength = &SJoinerError{"invalid length"}
)
