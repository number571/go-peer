package config

import "errors"

const (
	errPrefix = "cmd/hidden_lake/adapters/common/consumer/internal/config = "
)

var (
	ErrLoadLogging        = errors.New(errPrefix + "load logging")
	ErrInvalidConfig      = errors.New(errPrefix + "invalid config")
	ErrInitConfig         = errors.New(errPrefix + "init config")
	ErrDeserializeConfig  = errors.New(errPrefix + "deserialize config")
	ErrReadConfig         = errors.New(errPrefix + "read config")
	ErrConfigNotExist     = errors.New(errPrefix + "config not exist")
	ErrWriteConfig        = errors.New(errPrefix + "write config")
	ErrConfigAlreadyExist = errors.New(errPrefix + "config already exist")
)