package settings

const (
	CServiceName  = "HLT"
	CTitlePattern = "go-peer/hidden-lake-traffic"
)

const (
	CPathDB  = "hlt.db"
	CPathCFG = "hlt.cfg"
)

const (
	CDefaultCapMessages = (1 << 10)
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
