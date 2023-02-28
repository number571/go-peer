package config

type IConfig interface {
	Logging() ILogging
	Network() string
	Address() string
	Connection() string
}

type ILogging interface {
	HasInfo() bool
	HasWarn() bool
	HasErro() bool
}
