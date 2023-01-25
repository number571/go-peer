package config

type IConfig interface {
	Network() string
	Address() string
	Connection() string
}
