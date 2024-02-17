package settings

import (
	"time"
)

const (
	CServiceName     = "HLS"
	CServiceFullName = "hidden-lake-service"
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
	CHeaderResponseMode = "Hl-Service-Response-Mode"
)

const (
	CHeaderResponseModeON  = "on" // default
	CHeaderResponseModeOFF = "off"
)

const (
	CRequestQueueCapacity = (1 << 10) // 1024 request_ids ~= 44KiB
)

const (
	CRetryEnqueue = 2
	CFetchTimeout = time.Minute
)

const (
	CDefaultMessageSize     = (8 << 10) // 8KiB
	CDefaultWorkSize        = 22        // bits
	CDefaultKeySize         = 4096      // bits
	CDefaultQueuePeriod     = 5000      // 5 seconds
	CDefaultQueueRandPeriod = 0         // 0 seconds
	CDefaultLimitVoidSize   = (4 << 10) // 4KiB
	CDefaultF2FDisabled     = false     // friend-to-friend
	CDefaultNetworkKey      = "default"
)

const (
	CDefaultTCPAddress  = "127.0.0.1:9571"
	CDefaultHTTPAddress = "127.0.0.1:9572"
)

const (
	CQueueCapacity     = (1 << 8) // 256 messages ~= 2MiB
	CQueuePoolCapacity = (1 << 5) //  32 messages ~= 256KiB
)

const (
	CNetworkQueueCapacity = (2 << 10) // 2048 hashes ~= 64KiB
	CNetworkMaxConns      = (1 << 8)  // 256 conns
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
	CHandleNetworkPubKeyPath  = "/api/network/pubkey"
)
