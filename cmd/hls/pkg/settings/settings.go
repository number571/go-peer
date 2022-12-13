package settings

const (
	CContentType  = "application/json"
	CHeaderPubKey = "Sender-Public-Key"
)

const (
	CAKeySize = 4096
)

const (
	CHandleIndex          = "/api/index"
	CHandleConfigConnects = "/api/config/connects"
	CHandleConfigFriends  = "/api/config/friends"
	CHandleNetworkOnline  = "/api/network/online"
	CHandleNetworkPush    = "/api/network/push"
	CHandleNodeKey        = "/api/node/key"
)

const (
	CErrorNone = iota + 1
	CErrorMethod
	CErrorDecode
	CErrorPubKey
	CErrorPrivKey
	CErrorMessage
	CErrorResponse
	CErrorBroadcast
	CErrorExist
	CErrorNotExist
	CErrorAction
	CErrorValue
	CErrorOpen
	CErrorRead
	CErrorWrite
	CErrorSize
	CErrorUnauth
)
