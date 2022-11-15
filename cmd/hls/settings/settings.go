package settings

import (
	"time"

	"github.com/number571/go-peer/settings"
)

const (
	CTitlePattern = settings.CGopeerPrefix + "hidden-lake-service"
	CContentType  = "application/json"
	CHeaderPubKey = "Sender-Public-Key"
)

const (
	CHandleRoot     = "/api"
	CHandleConnects = "/api/config/connects"
	CHandleFriends  = "/api/config/friends"
	CHandleOnline   = "/api/network/online"
	CHandlePush     = "/api/network/push"
	CHandlePubKey   = "/api/node/pubkey"
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
	CErrorMessage
	CErrorResponse
	CErrorBroadcast
	CErrorExist
	CErrorNotExist
	CErrorAction
	CErrorValue
)
