package settings

const (
	CTitlePattern = "go-peer/hidden-lake-traffic"
	CContentType  = "application/json"
)

const (
	CHashesPath = "/hashes"
	CLoadPath   = "/load"
	CPushPath   = "/push"
)

const (
	CHashesTemplate = "%s" + CHashesPath
	CLoadTemplate   = "%s" + CLoadPath
	CPushTemplate   = "%s" + CPushPath
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
