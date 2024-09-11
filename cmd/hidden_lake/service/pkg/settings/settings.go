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
	CPathKey = "hls.key"
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
	CDefaultMessageSize     = (8 << 10) // 8192-B
	CDefaultKeySize         = (4 << 10) // 4096-b
	CDefaultFetchTimeout    = 60_000    // 60 seconds
	CDefaultQueuePeriod     = 5_000     // 5 seconds
	CDefaultWorkSize        = 0         // bits
	CDefaultRandQueuePeriod = 0         // 0 seconds
	CDefaultRandMessageSize = 0         // 0 bytes
	CDefaultNetworkKey      = ""
)

const (
	CDefaultTCPAddress  = "127.0.0.1:9571"
	CDefaultHTTPAddress = "127.0.0.1:9572"
)

const (
	CQueueMainPoolCapacity = (1 << 8) // 256 messages ~= 2MiB
	CQueueRandPoolCapacity = (1 << 5) //  32 messages ~= 256KiB
)

const (
	CNetworkQueueCapacity = (2 << 10) // 2048 hashes ~= 64KiB
	CNetworkMaxConns      = (1 << 8)  // 256 conns
	CNetworkWriteTimeout  = 5 * time.Second
	CNetworkReadTimeout   = 5 * time.Second
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
