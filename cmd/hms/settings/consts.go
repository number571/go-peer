package settings

import "github.com/number571/go-peer/settings"

const (
	CTitlePattern = settings.CGopeerPrefix + "hidden-message-service"
	CContentType  = "application/json"
)

const (
	CSizePath = "/size"
	CLoadPath = "/load"
	CPushSize = "/push"
)

const (
	CSizeTemplate = "%s" + CSizePath
	CLoadTemplate = "%s" + CLoadPath
	CPushTemplate = "%s" + CPushSize
)

const (
	CSizeWork = 20        // bits
	CSizePack = (8 << 20) // 8MiB
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
