package settings

const (
	CServiceName  = "HLT"
	CTitlePattern = "go-peer/hidden-lake-traffic"
)

const (
	CPathYML = "hlt.yml"
	CPathDB  = "hlt.db"
)

const (
	CDefaultMessagesCapacity = (2 << 10)
)

const (
	CHandleIndexPath   = "/api/index"
	CHandleHashesPath  = "/api/hashes"
	CHandleMessagePath = "/api/message"
)

const (
	CDefaultTCPAddress  = "127.0.0.1:9581"
	CDefaultHTTPAddress = "127.0.0.1:9582"
)

const (
	CDefaultConnectionAddress = "127.0.0.1:9571"
)
