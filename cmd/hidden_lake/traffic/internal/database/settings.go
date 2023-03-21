package database

var (
	_ ISettings = &sSettings{}
)

const (
	cPath          = "database.db"
	cWorkSize      = 10
	cMessageSize   = (1 << 20)
	cLimitMessages = (1 << 10)
)

type SSettings sSettings
type sSettings struct {
	FPath          string
	FLimitMessages uint64
	FMessageSize   uint64
	FWorkSize      uint64
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FPath:          pSett.FPath,
		FWorkSize:      pSett.FWorkSize,
		FMessageSize:   pSett.FMessageSize,
		FLimitMessages: pSett.FLimitMessages,
	}).useDefaultValues()
}

func (p *sSettings) useDefaultValues() ISettings {
	if p.FPath == "" {
		p.FPath = cPath
	}
	if p.FWorkSize == 0 {
		p.FWorkSize = cWorkSize
	}
	if p.FMessageSize == 0 {
		p.FMessageSize = cMessageSize
	}
	if p.FLimitMessages == 0 {
		p.FLimitMessages = cLimitMessages
	}
	return p
}

func (p *sSettings) GetPath() string {
	return p.FPath
}

func (s *sSettings) GetLimitMessages() uint64 {
	return s.FLimitMessages
}

func (p *sSettings) GetMessageSize() uint64 {
	return p.FMessageSize
}

func (p *sSettings) GetWorkSize() uint64 {
	return p.FWorkSize
}
