package settings

const (
	CServiceName     = "HLT"
	CServiceFullName = "hidden-lake-traffic"
)

const (
	CPathYML = "hlt.yml"
	CPathDB  = "hlt.db"
)

const (
	CHandleIndexPath          = "/api/index"
	CHandleStoragePointerPath = "/api/storage/pointer"
	CHandleStorageHashesPath  = "/api/storage/hashes"
	CHandleNetworkMessagePath = "/api/network/message"
	CHandleConfigSettings     = "/api/config/settings"
)

const (
	CDefaultMessagesCapacity = 0
	CDefaultDatabaseEnabled  = false
)

const (
	CDefaultTCPAddress  = "127.0.0.1:9581"
	CDefaultHTTPAddress = "127.0.0.1:9582"
)

const (
	CDefaultConnectionAddress = "127.0.0.1:9571"
)
