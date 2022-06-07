package settings

const (
	CTitlePattern = "hidden-message-service"
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
	CErrorNone = iota + 1
	CErrorMethod
	CErrorDecode
	CErrorLoad
	CErrorPush
	CErrorMessage
	CErrorPackSize
	CErrorWorkSize
)
