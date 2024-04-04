package config

import "errors"

var (
	ErrInvalidConfig      = errors.New("invalid config")
	ErrLoadLogging        = errors.New("load logging")
	ErrInitConfig         = errors.New("init config")
	ErrDeserializeConfig  = errors.New("deserialize config")
	ErrReadConfig         = errors.New("read config")
	ErrConfigNotExist     = errors.New("config not exist")
	ErrWriteConfig        = errors.New("write config")
	ErrConfigAlreadyExist = errors.New("config already exist")
)
