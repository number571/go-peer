package encoding

const (
	errPrefix = "pkg/encoding = "
)

type SEncodingError struct {
	str string
}

func (err *SEncodingError) Error() string {
	return errPrefix + err.str
}

var (
	ErrDeserialize = &SEncodingError{"deserialize bytes"}
)
