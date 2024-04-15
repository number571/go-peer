package adapted

const (
	errPrefix = "cmd/hidden_lake/adapters/chatingar/internal/consumer/internal/adapted = "
)

type SAdaptedError struct {
	str string
}

func (err *SAdaptedError) Error() string {
	return errPrefix + err.str
}

var (
	ErrBuildRequest      = &SAdaptedError{"build request"}
	ErrBadRequest        = &SAdaptedError{"bad request"}
	ErrBadStatusCode     = &SAdaptedError{"bad status code"}
	ErrDecodeCount       = &SAdaptedError{"decode count"}
	ErrCountLtNull       = &SAdaptedError{"count < 0"}
	ErrLimitPage         = &SAdaptedError{"limit page"}
	ErrDecodeMessages    = &SAdaptedError{"ldecode messages"}
	ErrLoadCountComments = &SAdaptedError{"load count comments"}
)
