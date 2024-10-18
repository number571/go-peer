package database

const (
	errPrefix = "pkg/storage/database = "
)

type SDatabaseError struct {
	str string
}

func (err *SDatabaseError) Error() string {
	return errPrefix + err.str
}

var (
	ErrOpenDB     = &SDatabaseError{"open database"}
	ErrSetValueDB = &SDatabaseError{"set value to database"}
	ErrGetValueDB = &SDatabaseError{"get value from database"}
	ErrDelValueDB = &SDatabaseError{"del value from database"}
	ErrNotFound   = &SDatabaseError{"value not found"}
	ErrCloseDB    = &SDatabaseError{"close database"}
)
