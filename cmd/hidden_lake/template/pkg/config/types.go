package config

import "github.com/number571/go-peer/cmd/hidden_lake/template/internal/config"

type IConfigSettings interface {
	config.IConfigSettings
}

type SConfigSettings struct {
	config.SConfigSettings
}
