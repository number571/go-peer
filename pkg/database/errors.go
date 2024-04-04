package database

const (
	errPrefix = "pkg/database = "
)

type SDatabaseError struct {
	str string
}

func (err *SDatabaseError) Error() string { return errPrefix + err.str }

var (
	ErrOpenDB               = &SDatabaseError{"open database"}
	ErrReadSalt             = &SDatabaseError{"read salt value"}
	ErrReadSaltHash         = &SDatabaseError{"read salt hash"}
	ErrPushSalt             = &SDatabaseError{"push salt value"}
	ErrPushSaltHash         = &SDatabaseError{"push salt hash"}
	ErrInvalidSaltHash      = &SDatabaseError{"invalid salt hash"}
	ErrSetValueDB           = &SDatabaseError{"set value to database"}
	ErrGetValueDB           = &SDatabaseError{"get value from database"}
	ErrDelValueDB           = &SDatabaseError{"del value from database"}
	ErrCloseDB              = &SDatabaseError{"close database"}
	ErrRecoverDB            = &SDatabaseError{"recover database"}
	ErrInvalidEncryptedSize = &SDatabaseError{"invalid encrypted size"}
	ErrInvalidDataHash      = &SDatabaseError{"invalid data hash"}
)
