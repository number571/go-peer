package config

const (
	errPrefix = "cmd/hidden_lake/helpers/traffic/internal/config = "
)

type SConfigError struct {
	str string
}

func (err *SConfigError) Error() string { return errPrefix + err.str }

var (
	ErrInvalidConfig      = &SConfigError{"invalid config"}
	ErrLoadLogging        = &SConfigError{"load logging"}
	ErrInitConfig         = &SConfigError{"init config"}
	ErrDeserializeConfig  = &SConfigError{"deserialize config"}
	ErrReadConfig         = &SConfigError{"read config"}
	ErrConfigNotExist     = &SConfigError{"config not exist"}
	ErrWriteConfig        = &SConfigError{"write config"}
	ErrConfigAlreadyExist = &SConfigError{"config already exist"}
)
