package config

type IConfig interface {
	Address() string
	Connection() string
}
