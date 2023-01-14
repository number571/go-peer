package settings

import (
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

const (
	CPathDB  = "hlt.db"
	CPathCFG = "hlt.cfg"
)

const (
	CLimitMessages = (1 << 10)
	CWorkSize      = hls_settings.CWorkSize
	CMessageSize   = hls_settings.CMessageSize
)
