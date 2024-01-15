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
	CHeaderSenderId  = "Hl-Messenger-Sender-Id"
	CHeaderRequestId = "Hl-Messenger-Request-Id"
)

const (
	CRequestIDSize = 32 // string chars (ASCII bytes)
)

const (
	// fourth digits of PI
	CAuthSalt = "4811174502_8410270193_8521105559_6446229489_5493038196"

	// fifth digits of PI
	CCipherSalt = "4428810975_6659334461_2847564823_3786783165_2712019091"
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
