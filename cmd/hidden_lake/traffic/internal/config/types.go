package config

type IConfig interface {
	Logging() ILogging
	Network() string
	Address() string
	Connection() string
}

type ILogging interface {
	Info() bool
	Warn() bool
	Erro() bool
}
