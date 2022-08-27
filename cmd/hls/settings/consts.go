package settings

import (
	"time"

	"github.com/number571/go-peer/settings"
)

const (
	CTitlePattern = settings.CGopeerPrefix + "hidden-lake-service"
	CContentType  = "application/json"
)

const (
	CHandleOnline  = "/online"
	CHandleRequest = "/do/request"
)

const (
	CHeaderHLS = uint32(0x1750571)
	CAKeySize  = 4096
)

const (
	CPathCFG = "hls.cfg"
	CPathSTG = "hls.stg"
	CPathDB  = "hls.db"
)

const (
	CRetryEnqueue = 2
	CWaitTime     = time.Minute
	CWorkSize     = 20        // bits
	CMessageSize  = (8 << 20) // 8MiB
)

const (
	CNetworkCapacity    = (4 << 10) // hashes
	CNetworkRetry       = 10        // retryNum for get message
	CNetworkMaxMessages = 20        // for one client
	CNetworkMaxConns    = 10
	CNetworkWaitTime    = 10 * time.Second
)

const (
	CQueueCapacity     = (1 << 8) // 2n messages in queue
	CQueuePullCapacity = (1 << 6) // generated fake messages
	CQueueDuration     = 5 * time.Second
)

const (
	CErrorNone = iota + 1
	CErrorMethod
	CErrorDecode
	CErrorPubKey
	CErrorResponse
)
