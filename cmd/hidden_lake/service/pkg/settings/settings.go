package settings

import "time"

const (
	CServiceName  = "HLS"
	CTitlePattern = "go-peer/hidden-lake-service"
)

const (
	CHeaderHLS = uint32(0x1750571)
)

const (
	CPathCFG = "hls.cfg"
	CPathDB  = "hls.db"
)

const (
	CHeaderPubKey  = "Service-Public-Key"
	CHeaderMsgHash = "Service-Message-Hash"
	CNetworkMask   = 0x676F2D7065657201
)

const (
	CRetryEnqueue = 2
	CWaitTime     = time.Minute
)

const (
	CQueueCapacity     = (1 << 6) // messages in queue
	CQueuePullCapacity = (1 << 5) // generated fake messages
	CQueueDuration     = 5 * time.Second
)

const (
	CAKeySize        = 4096
	CNetworkCapacity = (1 << 10) // hashes
	CNetworkMaxConns = 10
	CWorkSize        = 20        // bits
	CMessageSize     = (1 << 20) // 1MiB
	CMaxVoidSize     = (1 << 20) // 1MiB
	CNetworkWaitTime = 10 * time.Second
)

const (
	CHandleIndexPath          = "/api/index"
	CHandleConfigConnectsPath = "/api/config/connects"
	CHandleConfigFriendsPath  = "/api/config/friends"
	CHandleNetworkOnlinePath  = "/api/network/online"
	CHandleNetworkRequestPath = "/api/network/request"
	CHandleNodeKeyPath        = "/api/node/key"
)

const (
	CHandleIndexTemplate          = "%s" + CHandleIndexPath
	CHandleConfigConnectsTemplate = "%s" + CHandleConfigConnectsPath
	CHandleConfigFriendsTemplate  = "%s" + CHandleConfigFriendsPath
	CHandleNetworkOnlineTemplate  = "%s" + CHandleNetworkOnlinePath
	CHandleNetworkRequestTemplate = "%s" + CHandleNetworkRequestPath
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
