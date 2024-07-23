package database

const (
	errPrefix = "pkg/database = "
)

type SDatabaseError struct {
	str string
}

func (err *SDatabaseError) Error() string {
	return errPrefix + err.str
}

var (
	ErrOpenDB   = &SDatabaseError{"open database"}
	ErrCloseDB  = &SDatabaseError{"close database"}
	ErrNotFound = &SDatabaseError{"value not found"}
	ErrSetValue = &SDatabaseError{"set value"}
	ErrGetValue = &SDatabaseError{"get value"}
	ErrDelValue = &SDatabaseError{"del value"}
)
