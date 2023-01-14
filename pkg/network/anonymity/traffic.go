package anonymity

var (
	_ ITraffic = &sTraffic{}
)

type sTraffic struct {
	FDownload IDownloadF
	FUpload   IUploadF
}

func NewTraffic(d IDownloadF, u IUploadF) ITraffic {
	return &sTraffic{
		FDownload: d,
		FUpload:   u,
	}
}

func (t *sTraffic) Download() IDownloadF {
	return t.FDownload
}

func (t *sTraffic) Upload() IUploadF {
	return t.FUpload
}
