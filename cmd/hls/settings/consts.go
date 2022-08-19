package settings

const (
	CTitlePattern = "hidden-lake-service"
	CContentType  = "application/json"
)

const (
	CHeaderHLS = uint32(0x1750571)
	CAKeySize  = 4096
)

const (
	CWorkSize  = 20        // bits
	CWaitTime  = 60        // seconds
	CPackSize  = (8 << 20) // 8MiB
	CMaxConns  = 10
	CMaxMsgs   = 20   // for one client
	CQueueSize = 200  // 2n messages in queue
	CQueuePull = 50   // generated fake messages
	CQueueTime = 5000 // milliseconds
)

const (
	CErrorNone = iota + 1
	CErrorMethod
	CErrorDecode
	CErrorPubKey
	CErrorResponse
)
