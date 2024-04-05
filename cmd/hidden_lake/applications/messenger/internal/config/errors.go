package config

const (
	errPrefix = "cmd/hidden_lake/applications/messenger/internal/config = "
)

type SConfigError struct {
	str string
}

func (err *SConfigError) Error() string { return errPrefix + err.str }

var (
	ErrInvalidConfig      = &SConfigError{"invalid config"}
	ErrLoadLogging        = &SConfigError{"load logging"}
	ErrLoadLanguage       = &SConfigError{"load language"}
	ErrToLanguage         = &SConfigError{"to language"}
	ErrInitConfig         = &SConfigError{"init config"}
	ErrDeserializeConfig  = &SConfigError{"deserialize config"}
	ErrReadConfig         = &SConfigError{"read config"}
	ErrConfigNotExist     = &SConfigError{"config not exist"}
	ErrWriteConfig        = &SConfigError{"write config"}
	ErrConfigAlreadyExist = &SConfigError{"config already exist"}
)
