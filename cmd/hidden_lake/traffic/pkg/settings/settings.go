package settings

import "time"

const (
	CServiceName  = "HLT"
	CTitlePattern = "go-peer/hidden-lake-traffic"
)

const (
	CPathDB  = "hlt.db"
	CPathCFG = "hlt.cfg"
)

const (
	CCapacity = (1 << 10)
	CTimeWait = 10 * time.Second
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
