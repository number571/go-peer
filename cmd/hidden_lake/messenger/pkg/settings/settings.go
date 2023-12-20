package settings

const (
	CServiceName  = "HLM"
	CTitlePattern = "go-peer/hidden-lake-messenger"
)

const (
	CPathYML = "hlm.yml"
	CPathDB  = "hlm.db"
)

const (
	CStaticPath = "/static/"
	CPushPath   = "/push"
)

const (
	CIamAliasName    = "__iam__"
	CDefaultLanguage = "ENG"
)

const (
	CHeaderSenderId = "Hl-Messenger-Sender-Id"
)

const (
	CDefaultInterfaceAddress = "127.0.0.1:9591"
	CDefaultIncomingAddress  = "127.0.0.1:9592"
)

const (
	CDefaultMessagesCapacity = (2 << 10) // count
)

const (
	CHandleIndexPath          = "/"
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
