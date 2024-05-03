package response

const (
	errPrefix = "cmd/hidden_lake/service/pkg/response = "
)

type SResponseError struct {
	str string
}

func (err *SResponseError) Error() string {
	return errPrefix + err.str
}

var (
	ErrLoadBytesJoiner = &SResponseError{"load bytes joiner"}
	ErrDecodeResponse  = &SResponseError{"decode response"}
	ErrUnknownType     = &SResponseError{"unknown type"}
)
