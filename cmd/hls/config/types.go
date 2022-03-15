package config

type IConfig interface {
	Address() string
	Connections() []string
	GetService(string) (string, bool)
}
