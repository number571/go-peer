package settings

import (
	"time"
)

const (
	CTitlePattern = "go-peer/hidden-lake-service"
	CContentType  = "application/json"
	CHeaderPubKey = "Sender-Public-Key"
)

const (
	CHandleConfigConnects = "/api/config/connects"
	CHandleConfigFriends  = "/api/config/friends"
	CHandleNetworkOnline  = "/api/network/online"
	CHandleNetworkPush    = "/api/network/push"
	CHandleNodePubkey     = "/api/node/pubkey"
)

const (
	CHeaderHLS = uint32(0x1750571)
	CAKeySize  = 4096
)

const (
	CPathCFG = "hls.cfg"
	CPathDB  = "hls.db"
)

const (
	CRetryEnqueue = 2
	CWaitTime     = time.Minute
	CWorkSize     = 20        // bits
	CMessageSize  = (2 << 20) // 2MiB
)

const (
	CNetworkCapacity = (1 << 10) // hashes
	CNetworkMaxConns = 10
	CNetworkWaitTime = 10 * time.Second
)

const (
	CQueueCapacity     = (1 << 6) // messages in queue
	CQueuePullCapacity = (1 << 5) // generated fake messages
	CQueueDuration     = 5 * time.Second
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
)
