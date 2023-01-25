package settings

import (
	"time"
)

const (
	CHeaderHLS = uint32(0x1750571)
)

const (
	CPathCFG = "hls.cfg"
	CPathDB  = "hls.db"
)

const (
	CRetryEnqueue = 2
	CWaitTime     = time.Minute
)

const (
	CNetworkCapacity = (1 << 10) // hashes
	CNetworkMaxConns = 10
)

const (
	CQueueCapacity     = (1 << 6) // messages in queue
	CQueuePullCapacity = (1 << 5) // generated fake messages
	CQueueDuration     = 5 * time.Second
)
