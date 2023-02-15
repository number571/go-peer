package settings

import (
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

const (
	CServiceName  = "HLT"
	CTitlePattern = "go-peer/hidden-lake-traffic"
)

const (
	CPathDB  = "hlt.db"
	CPathCFG = "hlt.cfg"
)

const (
	CLimitMessages   = (1 << 10)
	CWorkSize        = hls_settings.CWorkSize
	CMessageSize     = hls_settings.CMessageSize
	CNetworkWaitTime = hls_settings.CNetworkWaitTime
	CNetworkMask     = hls_settings.CNetworkMask
)

const (
	CHandleIndexPath   = "/api/index"
	CHandleHashesPath  = "/api/hashes"
	CHandleMessagePath = "/api/message"
)

const (
	CHandleIndexTemplate   = "%s" + CHandleIndexPath
	CHandleHashesTemplate  = "%s" + CHandleHashesPath
	CHandleMessageTemplate = "%s" + CHandleMessagePath
)

const (
	CErrorNone = iota + 1
	CErrorMethod
	CErrorDecode
	CErrorLoad
	CErrorPush
	CErrorMessage
	CErrorPackSize
	CErrorWorkSize
)
