package language

const (
	errPrefix = "internal/language = "
)

type SLanguageError struct {
	str string
}

func (err *SLanguageError) Error() string {
	return errPrefix + err.str
}

var (
	ErrUnknownLanguage = &SLanguageError{"unknown language"}
)
