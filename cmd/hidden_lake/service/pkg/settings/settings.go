package settings

const (
	CTitlePattern = "go-peer/hidden-lake-service"
	CContentType  = "application/json"
	CHeaderPubKey = "Sender-Public-Key"
)

const (
	CAKeySize    = 4096
	CWorkSize    = 20        // bits
	CMessageSize = (2 << 20) // 2MiB
)

const (
	CHandleIndexPath          = "/api/index"
	CHandleConfigConnectsPath = "/api/config/connects"
	CHandleConfigFriendsPath  = "/api/config/friends"
	CHandleNetworkOnlinePath  = "/api/network/online"
	CHandleNetworkPushPath    = "/api/network/push"
	CHandleNodeKeyPath        = "/api/node/key"
)

const (
	CHandleIndexTemplate          = "%s" + CHandleIndexPath
	CHandleConfigConnectsTemplate = "%s" + CHandleConfigConnectsPath
	CHandleConfigFriendsTemplate  = "%s" + CHandleConfigFriendsPath
	CHandleNetworkOnlineTemplate  = "%s" + CHandleNetworkOnlinePath
	CHandleNetworkPushTemplate    = "%s" + CHandleNetworkPushPath
	CHandleNodeKeyTemplate        = "%s" + CHandleNodeKeyPath
)

const (
	CErrorNone = iota + 1
	CErrorMethod
	CErrorDecode
	CErrorPubKey
	CErrorPrivKey
	CErrorMessage
	CErrorResponse
	CErrorBroadcast
	CErrorExist
	CErrorNotExist
	CErrorAction
	CErrorValue
	CErrorOpen
	CErrorRead
	CErrorWrite
	CErrorSize
	CErrorUnauth
)
