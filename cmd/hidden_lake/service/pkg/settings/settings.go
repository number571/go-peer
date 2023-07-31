package settings

import (
	"time"

	"github.com/number571/go-peer/pkg/network/anonymity/logbuilder"
)

const (
	CServiceName  = "HLS"
	CTitlePattern = "go-peer/hidden-lake-service"
)

const (
	CServiceMask = uint32(0x1750571)
	CNetworkMask = uint64(0x676F2D7065657201)
)

const (
	CPathCFG = "hls.cfg"
	CPathDB  = "hls.db"
)

const (
	CHeaderPublicKey   = "Service-Public-Key"
	CHeaderMessageHash = "Service-Message-Hash"
	CHeaderOffResponse = "Service-Off-Response"
)

const (
	CRetryEnqueue = 2
	CFetchTimeout = time.Minute
)

const (
	CDefaultMessageSize   = (4 << 20) // 4KiB
	CDefaultWorkSize      = 20        // bits
	CDefaultKeySize       = 4096      // bits
	CDefaultQueuePeriod   = 5000      // 5seconds
	CDefaultLimitVoidSize = (4 << 10) // 4KiB
)

const (
	CDefaultTCPAddress        = "127.0.0.1:9571"
	CDefaultHTTPAddress       = "127.0.0.1:9572"
	CDefaultServiceHLMAddress = "127.0.0.1:9592"
)

const (
	CQueueCapacity     = (1 << 6) // messages in queue
	CQueuePoolCapacity = (1 << 5) // generated fake messages
)

const (
	CNetworkCapacity = (1 << 10) // 1024 hashes
	CNetworkMaxConns = (1 << 6)  // 64 conns
)

const (
	CConnWaitReadDeadline = time.Hour
	CConnKeeperDuration   = 10 * time.Second
)

const (
	CHandleIndexPath          = "/api/index"
	CHandleConfigConnectsPath = "/api/config/connects"
	CHandleConfigFriendsPath  = "/api/config/friends"
	CHandleNetworkOnlinePath  = "/api/network/online"
	CHandleNetworkRequestPath = "/api/network/request"
	CHandleNetworkMessagePath = "/api/network/message"
	CHandleNodeKeyPath        = "/api/node/key"
)

const (
	CHandleIndexTemplate          = "%s" + CHandleIndexPath
	CHandleConfigConnectsTemplate = "%s" + CHandleConfigConnectsPath
	CHandleConfigFriendsTemplate  = "%s" + CHandleConfigFriendsPath
	CHandleNetworkOnlineTemplate  = "%s" + CHandleNetworkOnlinePath
	CHandleNetworkRequestTemplate = "%s" + CHandleNetworkRequestPath
	CHandleNetworkMessageTemplate = "%s" + CHandleNetworkMessagePath
	CHandleNodeKeyTemplate        = "%s" + CHandleNodeKeyPath
)

const (
	// INFO
	CLogInfoOffResponseFromService logbuilder.ILogType = "OFRSP"
	CLogInfoResponseFromService    logbuilder.ILogType = "RSPSR"

	// WARN
	CLogWarnRequestToService logbuilder.ILogType = "RQTSR"
	CLogWarnUndefinedService logbuilder.ILogType = "UNDSR"

	// ERRO
	CLogErroLoadRequestType  logbuilder.ILogType = "LDRQT"
	CLogErroProxyRequestType logbuilder.ILogType = "PXRQT"
)
