package config

type IConfig interface {
	GetLogging() ILogging
	GetNetwork() string
	GetAddress() string
	GetConnection() string
	GetConsumers() []string
}

type ILogging interface {
	HasInfo() bool
	HasWarn() bool
	HasErro() bool
}
