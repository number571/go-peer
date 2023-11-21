package settings

const (
	CServiceName  = "HLL"
	CTitlePattern = "go-peer/hidden-lake-loader"
)

const (
	CPathCFG = "hll.cfg"
)

const (
	CDefaultHTTPAddress = "127.0.0.1:9561"
)

const (
	CDefaultConsumerAddress = "127.0.0.1:9582"
	CDefaultProducerAddress = "127.0.0.2:9582"
)

const (
	CHandleIndexPath    = "/api/index"
	CHandleTransferPath = "/api/transfer"
)

const (
	CHandleIndexTemplate    = "%s" + CHandleIndexPath
	CHandleTransferTemplate = "%s" + CHandleTransferPath
)
