package client

const (
	errPrefix = "cmd/hidden_lake/helpers/template/pkg/client = "
)

type SClientError struct {
	str string
}

func (err *SClientError) Error() string {
	return errPrefix + err.str
}

var (
	ErrBadRequest     = &SClientError{"bad request"}
	ErrDecodeResponse = &SClientError{"decode response"}
	ErrInvalidTitle   = &SClientError{"invalid title"}
)
