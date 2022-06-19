package errors

var (
	_ error = &sError{}
)

type sError struct {
	fType uint64
	fMsg  string
}

func NewError(errType uint64, errMsg string) error {
	return &sError{
		fType: errType,
		fMsg:  errMsg,
	}
}

func IsError(err error, errType uint64) bool {
	if err == nil {
		return false
	}
	switch x := err.(type) {
	case *sError:
		return x.fType == errType
	default:
		return false
	}
}

func (err *sError) Error() string {
	return err.fMsg
}
