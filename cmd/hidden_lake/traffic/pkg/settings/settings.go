package settings

const (
	CTitlePattern = "go-peer/hidden-lake-traffic"
	CContentType  = "application/json"
)

const (
	CHandleIndexPath     = "/api/index"
	CHandleHashesPath    = "/api/hashes"
	CHandleMessagePath   = "/api/message"
	CHandleBroadcastPath = "/api/broadcast"
)

const (
	CHandleIndexTemplate     = "%s" + CHandleIndexPath
	CHandleHashesTemplate    = "%s" + CHandleHashesPath
	CHandleMessageTemplate   = "%s" + CHandleMessagePath
	CHandleBroadcastTemplate = "%s" + CHandleBroadcastPath
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
