package client

type sSettings struct {
	fWorkSize    uint64
	fMessageSize uint64
}

func NewSettings(workSize, msgSize uint64) ISettings {
	return &sSettings{
		fWorkSize:    workSize,
		fMessageSize: msgSize,
	}
}

func (s *sSettings) GetWorkSize() uint64 {
	return s.fWorkSize
}

func (s *sSettings) GetMessageSize() uint64 {
	return s.fMessageSize
}
