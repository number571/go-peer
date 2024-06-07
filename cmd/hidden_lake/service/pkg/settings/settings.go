package settings

import (
	"time"
)

const (
	CServiceName     = "HLS"
	CServiceFullName = "hidden-lake-service"
)

const (
	CNetworkMask = uint32(0x5f67705f) // bytes_prefix: _gp_
	CServiceMask = uint32(0x5f686c5f) // bytes_prefix: _hl_
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
	CRetryEnqueue   = 2
	CFetchTimeRatio = 10
)

const (
	CDefaultMessageSize     = (8 << 10) // 8KiB
	CDefaultWorkSize        = 22        // bits
	CDefaultKeySize         = 4096      // bits
	CDefaultQueuePeriod     = 5000      // 5 seconds
	CDefaultQueueRandPeriod = 0         // 0 seconds
	CDefaultLimitVoidSize   = (4 << 10) // 4KiB
	CDefaultF2FDisabled     = false     // friend-to-friend
	CDefaultQBTDisabled     = false     // queue-based-task
	CDefaultNetworkKey      = "default"
)

const (
	CDefaultTCPAddress  = "127.0.0.1:9571"
	CDefaultHTTPAddress = "127.0.0.1:9572"
)

const (
	CQueueMainCapacity = (1 << 8) // 256 messages ~= 2MiB
	CQueueVoidCapacity = (1 << 5) //  32 messages ~= 256KiB
)

const (
	CNetworkQueueCapacity = (2 << 10) // 2048 hashes ~= 64KiB
	CNetworkMaxConns      = (1 << 8)  // 256 conns
)

const (
	CConnKeeperDuration  = 10 * time.Second
	CConnDialTimeout     = 30 * time.Second
	CConnWaitReadTimeout = time.Hour
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
