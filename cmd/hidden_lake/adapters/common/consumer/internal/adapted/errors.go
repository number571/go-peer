package adapted

const (
	errPrefix = "cmd/hidden_lake/adapters/common/consumer/internal/adapted = "
)

type SAdaptedError struct {
	str string
}

func (err *SAdaptedError) Error() string {
	return errPrefix + err.str
}

var (
	ErrLoadCountService = &SAdaptedError{"load count from service"}
	ErrLoadCountDB      = &SAdaptedError{"load count from db"}
	ErrGetCount         = &SAdaptedError{"get count"}
	ErrSetNewCount      = &SAdaptedError{"set new count"}
	ErrInitCountKey     = &SAdaptedError{"init count key"}
	ErrParseCount       = &SAdaptedError{"parse count"}
	ErrInvalidResponse  = &SAdaptedError{"invalid response"}
	ErrReadResponse     = &SAdaptedError{"read response"}
	ErrBadRequest       = &SAdaptedError{"bad request"}
	ErrBuildRequest     = &SAdaptedError{"build request"}
	ErrDecodeMessage    = &SAdaptedError{"decode message"}
	ErrIncrementCount   = &SAdaptedError{"increment count"}
	ErrLoadMessage      = &SAdaptedError{"load message"}
)
