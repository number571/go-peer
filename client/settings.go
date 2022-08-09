package client

type sSettings struct {
	fWorkSize uint64
	fRandSize uint64
}

func NewSettings(workSize, randSize uint64) ISettings {
	return &sSettings{
		fWorkSize: workSize,
		fRandSize: randSize,
	}
}

func (s *sSettings) GetWorkSize() uint64 {
	return s.fWorkSize
}

func (s *sSettings) GetRandomSize() uint64 {
	return s.fRandSize
}
