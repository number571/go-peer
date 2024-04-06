package api

const (
	errPrefix = "internal/api = "
)

type SApiError struct {
	str string
}

func (err *SApiError) Error() string {
	return errPrefix + err.str
}

var (
	ErrBadStatusCode = &SApiError{"bad status code"}
	ErrReadResponse  = &SApiError{"read response"}
	ErrLoadResponse  = &SApiError{"load response"}
	ErrBadRequest    = &SApiError{"bad request"}
	ErrBuildRequest  = &SApiError{"build request"}
	ErrCopyBytes     = &SApiError{"copy bytes"}
)
