package config

type IConfig interface {
	GetAddress() IAddress
	GetConnection() IConnection
	GetStorageKey() string
}

type IAddress interface {
	GetInterface() string
	GetIncoming() string
}

type IConnection interface {
	GetService() string
	GetTraffic() string
}
