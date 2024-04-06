package state

const (
	errPrefix = "pkg/state = "
)

type SStateError struct {
	str string
}

func (err *SStateError) Error() string {
	return errPrefix + err.str
}

var (
	ErrStateEnabled  = &SStateError{"state already enabled"}
	ErrStateDisabled = &SStateError{"state already disabled"}
	ErrFuncEnable    = &SStateError{"enable state function"}
	ErrFuncDisable   = &SStateError{"disable state function"}
)
