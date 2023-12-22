package settings

import (
	"time"
)

const (
	CServiceName  = "HLS"
	CTitlePattern = "go-peer/hidden-lake-service"
)

const (
	CServiceMask = uint32(0x1750571)
	CNetworkMask = uint64(0x676F2D7065657201) // bytes_prefix: go-peer
)

const (
	CPathYML = "hls.yml"
	CPathDB  = "hls.db"
)

const (
	CHeaderPublicKey    = "Hl-Service-Public-Key"
	CHeaderRequestId    = "Hl-Service-Request-Id"
	CHeaderResponseMode = "Hl-Service-Response-Mode"
)

const (
	CHeaderResponseModeON  = "on" // default
	CHeaderResponseModeOFF = "off"
)

const (
	CHandleRequestQueueSize = (1 << 10) // 1024 values
	CHandleRequestIDSize    = 16        // bytes
)

const (
	CRetryEnqueue = 2
	CFetchTimeout = time.Minute
)

const (
	CDefaultMessageSize   = (8 << 10) // 8KiB
	CDefaultWorkSize      = 20        // bits
	CDefaultKeySize       = 4096      // bits
	CDefaultQueuePeriod   = 5000      // 5seconds
	CDefaultLimitVoidSize = (4 << 10) // 4KiB
)

const (
	CDefaultTCPAddress  = "127.0.0.1:9571"
	CDefaultHTTPAddress = "127.0.0.1:9572"
)

const (
	CDefaultServiceHLTAddress = "127.0.0.1:9582"
	CDefaultServiceHLMAddress = "127.0.0.1:9592"
)

const (
	CQueueCapacity     = (1 << 6) // messages in queue
	CQueuePoolCapacity = (1 << 5) // generated fake messages
)

const (
	CNetworkQueueSize = (1 << 10) // 1024 hashes
	CNetworkMaxConns  = (1 << 8)  // 256 conns
)

const (
	CConnWaitReadDeadline = time.Hour
	CConnKeeperDuration   = 10 * time.Second
)

const (
	CHandleIndexPath          = "/api/index"
	CHandleConfigSettingsPath = "/api/config/settings"
	CHandleConfigConnectsPath = "/api/config/connects"
	CHandleConfigFriendsPath  = "/api/config/friends"
	CHandleNetworkOnlinePath  = "/api/network/online"
	CHandleNetworkRequestPath = "/api/network/request"
	CHandleNetworkKeyPath     = "/api/network/key"
	CHandleNodeKeyPath        = "/api/node/key"
)

func GetConnDeadline(pQueuePeriod time.Duration) time.Duration {
	return pQueuePeriod / 2
}
