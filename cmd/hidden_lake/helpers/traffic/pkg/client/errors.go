package client

const (
	errPrefix = "cmd/hidden_lake/helpers/traffic/pkg/client = "
)

type SClientError struct {
	str string
}

func (err *SClientError) Error() string {
	return errPrefix + err.str
}

var (
	ErrBadRequest       = &SClientError{"bad request"}
	ErrDecodeResponse   = &SClientError{"decode response"}
	ErrInvalidResponse  = &SClientError{"invalid response"}
	ErrInvalidHexFormat = &SClientError{"invalid hex format"}
	ErrInvalidTitle     = &SClientError{"invalid title"}
	ErrDecodeMessage    = &SClientError{"decode message"}
)
