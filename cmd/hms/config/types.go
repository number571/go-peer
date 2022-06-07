package config

type IConfig interface {
	Address() string
	CleanCron() string
	Connections() []string
}
