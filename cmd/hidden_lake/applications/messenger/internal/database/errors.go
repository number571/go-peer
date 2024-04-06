package database

const (
	errPrefix = "cmd/hidden_lake/applications/messenger/internal/database = "
)

type SDatabaseError struct {
	str string
}

func (err *SDatabaseError) Error() string {
	return errPrefix + err.str
}

var (
	ErrLoadMessage    = &SDatabaseError{"load message"}
	ErrGetMessage     = &SDatabaseError{"get message"}
	ErrSetMessage     = &SDatabaseError{"set message"}
	ErrSetSizeMessage = &SDatabaseError{"set size message"}
	ErrCloseDB        = &SDatabaseError{"close db"}
	ErrEndGtSize      = &SDatabaseError{"end > size"}
	ErrStartGtEnd     = &SDatabaseError{"start > end"}
	ErrCreateDB       = &SDatabaseError{"create db"}
)
