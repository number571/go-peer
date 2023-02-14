package settings

import "time"

const (
	CTitlePattern  = "go-peer/hidden-lake-service"
	CContentType   = "application/json"
	CHeaderPubKey  = "Service-Public-Key"
	CHeaderMsgHash = "Service-Message-Hash"
	CNetworkMask   = 0x676F2D7065657201
)

const (
	CAKeySize        = 4096
	CWorkSize        = 20        // bits
	CMessageSize     = (1 << 20) // 1MiB
	CNetworkWaitTime = 10 * time.Second
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
