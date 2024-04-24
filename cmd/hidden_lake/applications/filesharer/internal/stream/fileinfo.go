package stream

var (
	_ IFileInfo = &sFileInfo{}
)

type sFileInfo struct {
	fName string
	fHash string
	fSize uint64
}

func NewFileInfo(pName, pHash string, pSize uint64) IFileInfo {
	return &sFileInfo{
		fName: pName,
		fHash: pHash,
		fSize: pSize,
	}
}

func (p *sFileInfo) GetName() string {
	return p.fName
}

func (p *sFileInfo) GetHash() string {
	return p.fHash
}

func (p *sFileInfo) GetSize() uint64 {
	return p.fSize
}
