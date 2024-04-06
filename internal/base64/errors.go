package base64

const (
	errPrefix = "internal/base64 = "
)

type SBase64Error struct {
	str string
}

func (err *SBase64Error) Error() string {
	return errPrefix + err.str
}

var (
	ErrBase64Size = &SBase64Error{"base64 size"}
)
