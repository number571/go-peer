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
	CDefaultTCPAddress   = "127.0.0.1:9581"
	CDefaultHTTPAddress  = "127.0.0.1:9582"
	CDefaultPPROFAddress = "127.0.0.1:9583"
)

const (
	CDefaultConnectionAddress = "127.0.0.1:9571"
)
