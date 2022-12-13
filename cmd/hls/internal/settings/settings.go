package settings

import (
	"time"
)

const (
	CTitlePattern = "go-peer/hidden-lake-service"
	CHeaderHLS    = uint32(0x1750571)
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
