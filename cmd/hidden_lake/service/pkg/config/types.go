package config

import "github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"

type IConfigSettings interface {
	config.IConfigSettings
	GetLimitMessageSizeBytes() uint64
}

type SConfigSettings struct {
	config.SConfigSettings
	FLimitMessageSizeBytes uint64 `json:"limit_message_size_bytes"`
}

func (p *SConfigSettings) GetLimitMessageSizeBytes() uint64 {
	return p.FLimitMessageSizeBytes
}
