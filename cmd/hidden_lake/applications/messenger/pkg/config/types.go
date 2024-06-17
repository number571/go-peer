package config

import "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/config"

type IConfigSettings interface {
	config.IConfigSettings
}

type SConfigSettings struct {
	config.SConfigSettings
}
