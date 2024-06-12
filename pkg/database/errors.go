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
	ErrOpenDB               = &SDatabaseError{"open database"}
	ErrGetSalt              = &SDatabaseError{"get salt value"}
	ErrReadSalt             = &SDatabaseError{"read salt value"}
	ErrReadRand             = &SDatabaseError{"read rand"}
	ErrGetHashRand          = &SDatabaseError{"get hash rand"}
	ErrReadHashRand         = &SDatabaseError{"read hash rand"}
	ErrPushSalt             = &SDatabaseError{"push salt value"}
	ErrPushRand             = &SDatabaseError{"push rand value"}
	ErrPushHashRand         = &SDatabaseError{"push hash rand"}
	ErrInvalidHash          = &SDatabaseError{"invalid hash"}
	ErrSetValueDB           = &SDatabaseError{"set value to database"}
	ErrGetValueDB           = &SDatabaseError{"get value from database"}
	ErrDelValueDB           = &SDatabaseError{"del value from database"}
	ErrGetNotFound          = &SDatabaseError{"get value not found"}
	ErrCloseDB              = &SDatabaseError{"close database"}
	ErrInvalidEncryptedSize = &SDatabaseError{"invalid encrypted size"}
	ErrInvalidDataHash      = &SDatabaseError{"invalid data hash"}
	ErrOpenBucket           = &SDatabaseError{"open bucket"}
)
