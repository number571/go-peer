package qb

const (
	errPrefix = "pkg/anonymity/qb = "
)

type SAnonymityError struct {
	str string
}

func (err *SAnonymityError) Error() string {
	return errPrefix + err.str
}

var (
	ErrSetHashIntoDB         = &SAnonymityError{"set hash into database"}
	ErrGetHashFromDB         = &SAnonymityError{"get hash from database"}
	ErrNilDB                 = &SAnonymityError{"database is nil"}
	ErrRetryLimit            = &SAnonymityError{"retry limit"}
	ErrEnqueueMessage        = &SAnonymityError{"enqueue message"}
	ErrUnknownType           = &SAnonymityError{"unknown type"}
	ErrLoadMessage           = &SAnonymityError{"load message"}
	ErrInvalidNetworkMask    = &SAnonymityError{"invalid network mask"}
	ErrStoreHashIntoDatabase = &SAnonymityError{"store hash into database"}
	ErrStoreHashWithProduce  = &SAnonymityError{"store hash with produce"}
	ErrActionIsNotFound      = &SAnonymityError{"action is not found"}
	ErrActionIsClosed        = &SAnonymityError{"action is closed"}
	ErrActionTimeout         = &SAnonymityError{"action timeout"}
	ErrEnqueuePayload        = &SAnonymityError{"enqueue payload"}
	ErrFetchResponse         = &SAnonymityError{"fetch response"}
	ErrRunning               = &SAnonymityError{"node running"}
	ErrProcessRun            = &SAnonymityError{"process run"}
	ErrHashAlreadyExist      = &SAnonymityError{"hash already exist"}
)
