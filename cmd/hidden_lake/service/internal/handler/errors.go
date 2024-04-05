package handler

const (
	errPrefix = "cmd/hidden_lake/service/internal/handler = "
)

type SHandlerError struct {
	str string
}

func (err *SHandlerError) Error() string { return errPrefix + err.str }

var (
	ErrBadRequest          = &SHandlerError{"bad request"}
	ErrBuildRequest        = &SHandlerError{"build request"}
	ErrUndefinedService    = &SHandlerError{"undefined service"}
	ErrLoadRequest         = &SHandlerError{"load request"}
	ErrUpdateFriends       = &SHandlerError{"update friends"}
	ErrInvalidResponseMode = &SHandlerError{"invalid response mode"}
)
