package settings

const (
	CServiceName  = "HLF"
	CTitlePattern = "hidden-lake-filesharer"
)

const (
	CPathYML       = "hlf.yml"
	CPathSTG       = "hlf.stg"
	CPathLoadedSTG = "hlf.stg/loaded"
)

const (
	CPageOffset = 10
	CChunkSize  = 3000 // bytes
)

// TODO: delete
// 4752 = message limit
// 4105 = response size (3000 bytes = 4000 base64 bytes)
// 105 = header bytes

const (
	CDefaultInterfaceAddress = "127.0.0.1:9541"
	CDefaultIncomingAddress  = "127.0.0.1:9542"
)

const (
	CDefaultLanguage = "ENG"
)

const (
	CHandleIndexPath          = "/"
	CHandleAboutPath          = "/about"
	CHandleFaviconPath        = "/favicon.ico"
	CHandleSettingsPath       = "/settings"
	CHandleFriendsPath        = "/friends"
	CHandleFriendsStoragePath = "/friends/storage"
	CStaticPath               = "/static/"
)

const (
	CListPath = "/list"
	CLoadPath = "/load"
)
