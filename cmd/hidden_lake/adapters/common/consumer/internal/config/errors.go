package config

import "errors"

var (
	ErrLoadLogging        = errors.New("load logging")
	ErrInvalidConfig      = errors.New("invalid config")
	ErrInitConfig         = errors.New("init config")
	ErrDeserializeConfig  = errors.New("deserialize config")
	ErrReadConfig         = errors.New("read config")
	ErrConfigNotExist     = errors.New("config not exist")
	ErrWriteConfig        = errors.New("write config")
	ErrConfigAlreadyExist = errors.New("config already exist")
)
