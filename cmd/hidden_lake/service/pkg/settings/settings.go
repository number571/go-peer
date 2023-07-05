package settings

import (
	"time"

	anon_logger "github.com/number571/go-peer/pkg/network/anonymity/logger"
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
	CAKeySize    = 4096      // bits
	CWorkSize    = 20        // bits
	CMessageSize = (1 << 20) // 1MiB
)

const (
	CQueueCapacity     = (1 << 6) // messages in queue
	CQueuePoolCapacity = (1 << 5) // generated fake messages
	CQueueDuration     = 5 * time.Second
)

const (
	CNetworkWriteTimeout = time.Minute
	CNetworkCapacity     = (1 << 10) // hashes
	CNetworkMaxConns     = (1 << 6)  // 64
)

const (
	CConnLimitVoidSize    = (1 << 20) // 1MiB
	CConnWaitReadDeadline = time.Hour
	CConnReadDeadline     = time.Minute
	CConnWriteDeadline    = time.Minute
)

const (
	CConnKeeperDuration = time.Minute
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
	// Base
	CLogWarnOffResponseFromService anon_logger.ILogType = "OFRSP"
	CLogWarnResponseFromService    anon_logger.ILogType = "RSPSR"

	// WARN
	CLogWarnRequestToService anon_logger.ILogType = "RQTSR"
	CLogWarnUndefinedService anon_logger.ILogType = "UNDSR"

	// ERRO
	CLogErroLoadRequestType  anon_logger.ILogType = "LDRQT"
	CLogErroProxyRequestType anon_logger.ILogType = "PXRQT"
)
