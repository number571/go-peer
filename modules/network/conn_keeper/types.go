package conn_keeper

import "time"

type IConnKeeper interface {
	Settings() ISettings

	Run() error
	Close() error
}

type ISettings interface {
	GetConnections() []string
	GetDuration() time.Duration
}
