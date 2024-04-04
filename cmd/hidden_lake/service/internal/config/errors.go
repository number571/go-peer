package config

import "errors"

const (
	errPrefix = "cmd/hidden_lake/service/internal/config = "
)

var (
	ErrNotSupportedKeySize = errors.New(errPrefix + "not supported key size")
	ErrInvalidPublicKey    = errors.New(errPrefix + "invalid public key")
	ErrDuplicatePublicKey  = errors.New(errPrefix + "duplicate public key")
	ErrLoadLogging         = errors.New(errPrefix + "load logging")
	ErrLoadPublicKey       = errors.New(errPrefix + "load public key")
	ErrLoadConfigSettings  = errors.New(errPrefix + "load config settings")
	ErrLoadConfig          = errors.New(errPrefix + "load config")
	ErrInitConfig          = errors.New(errPrefix + "init config")
	ErrDeserializeConfig   = errors.New(errPrefix + "deserialize config")
	ErrReadConfig          = errors.New(errPrefix + "read config")
	ErrConfigNotFound      = errors.New(errPrefix + "config not found")
	ErrWriteConfig         = errors.New(errPrefix + "write config")
	ErrConfigAlreadyExist  = errors.New(errPrefix + "config already exist")
	ErrBuildConfig         = errors.New(errPrefix + "build config")
)
