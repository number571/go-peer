package settings

const (
	CServiceName     = "HLR"
	CServiceFullName = "hidden-lake-remoter"
)

const (
	CPathYML        = "hlr.yml"
	CHeaderPassword = "Hl-Remoter-Password" // nolint: gosec
)

const (
	CExecPath      = "/exec"
	CExecSeparator = "[@remoter-separator]"
)

const (
	CDefaultIncomingAddress = "127.0.0.1:9532"
	CDefaultExecTimeout     = 5_000 // 5s
)
