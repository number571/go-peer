package settings

const (
	CTitlePattern = "hidden-lake-service"
	CContentType  = "application/json"
)

const (
	CHeaderHLS = uint32(0x1750571)
	CSizeRoute = 3
	CAKeySize  = 4096
)

const (
	CErrorNone = iota + 1
	CErrorMethod
	CErrorDecode
	CErrorPubKey
	CErrorResponse
)
