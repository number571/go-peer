package settings

const (
	CServiceName  = "HLM"
	CTitlePattern = "hidden-lake-messenger"
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
	CHeaderPseudonym = "Hl-Messenger-Pseudonym"
	CHeaderRequestId = "Hl-Messenger-Request-Id"
)

const (
	CPseudonymSize = 32 // x >= 1 && x <= CPseudonymSize
	CRequestIDSize = 44 // string chars (ASCII bytes)
)

const (
	CDefaultInterfaceAddress = "127.0.0.1:9591"
	CDefaultIncomingAddress  = "127.0.0.1:9592"
)

const (
	CDefaultShare            = false
	CDefaultLanguage         = "ENG"
	CDefaultMessagesCapacity = (2 << 10) // count
)

const (
	CHandleIndexPath         = "/"
	CHandleAboutPath         = "/about"
	CHandleFaviconPath       = "/favicon.ico"
	CHandleSettingsPath      = "/settings"
	CHandleFriendsPath       = "/friends"
	CHandleFriendsChatPath   = "/friends/chat"
	CHandleFriendsUploadPath = "/friends/upload"
	CHandleFriendsChatWSPath = "/friends/chat/ws"
)

const (
	CIsText = 0x01
	CIsFile = 0x02
)
