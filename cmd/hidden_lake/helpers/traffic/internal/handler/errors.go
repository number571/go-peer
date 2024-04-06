package handler

const (
	errPrefix = "cmd/hidden_lake/helpers/traffic/internal/handler = "
)

type SHandlerError struct {
	str string
}

func (err *SHandlerError) Error() string {
	return errPrefix + err.str
}

var (
	ErrLoadMessage   = &SHandlerError{"load message"}
	ErrDatabaseNull  = &SHandlerError{"database null"}
	ErrPushMessageDB = &SHandlerError{"push message db"}
)
