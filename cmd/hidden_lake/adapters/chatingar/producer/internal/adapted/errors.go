package adapted

const (
	errPrefix = "cmd/hidden_lake/adapters/chatingar/producer/internal/adapted = "
)

type SAdaptedError struct {
	str string
}

func (err *SAdaptedError) Error() string {
	return errPrefix + err.str
}

var (
	ErrBadStatusCode = &SAdaptedError{"bad status code"}
	ErrBadRequest    = &SAdaptedError{"bad request"}
	ErrBuildRequest  = &SAdaptedError{"build request"}
)
