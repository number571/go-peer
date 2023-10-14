package settings

const (
	CServiceName  = "HLM"
	CTitlePattern = "go-peer/hidden-lake-messenger"
)

const (
	CPathCFG = "hlm.cfg"
	CPathDB  = "hlm.db"
)

const (
	CStaticPath = "/static/"
	CPushPath   = "/push"
)

const (
	CIamAliasName      = "__iam__"
	CDefaultLanguage   = "ENG"
	CDefaultStorageKey = "_"
)

const (
	CDefaultInterfaceAddress = "127.0.0.1:9591"
	CDefaultIncomingAddress  = "127.0.0.1:9592"
)

const (
	CMinEntropy         = 64        // bits
	CWorkForKeys        = 20        // bits
	CDefaultCapMessages = (2 << 10) // count
)

const (
	CHandleIndexPath          = "/"
	CHandleSignOutPath        = "/sign/out"
	CHandleSignInPath         = "/sign/in"
	CHandleSignUpPath         = "/sign/up"
	CHandleAboutPath          = "/about"
	CHandleFaviconPath        = "/favicon.ico"
	CHandleSettingsPath       = "/settings"
	CHandleQRPublicKeyKeyPath = "/qr/public_key"
	CHandleFriendsPath        = "/friends"
	CHandleFriendsChatPath    = "/friends/chat"
	CHandleFriendsUploadPath  = "/friends/upload"
	CHandleFriendsChatWSPath  = "/friends/chat/ws"
)

const (
	CIsText = 0x01
	CIsFile = 0x02
)
