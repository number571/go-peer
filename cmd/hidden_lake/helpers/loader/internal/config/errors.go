package config

import "errors"

const (
	errPrefix = "cmd/hidden_lake/helpers/loader/internal/config = "
)

var (
	ErrInvalidConfig      = errors.New(errPrefix + "invalid config")
	ErrLoadLogging        = errors.New(errPrefix + "load logging")
	ErrInitConfig         = errors.New(errPrefix + "init config")
	ErrDeserializeConfig  = errors.New(errPrefix + "deserialize config")
	ErrReadConfig         = errors.New(errPrefix + "read config")
	ErrConfigNotExist     = errors.New(errPrefix + "config not exist")
	ErrWriteConfig        = errors.New(errPrefix + "write config")
	ErrConfigAlreadyExist = errors.New(errPrefix + "config already exist")
)
