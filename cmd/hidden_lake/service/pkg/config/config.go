package config

var (
	_ IConfigSettings = &SConfigSettings{}
)

type SConfigSettings struct {
	FMessageSizeBytes   uint64 `json:"message_size_bytes"`
	FWorkSizeBits       uint64 `json:"work_size_bits"`
	FQueuePeriodMS      uint64 `json:"queue_period_ms"`
	FKeySizeBits        uint64 `json:"key_size_bits"`
	FLimitVoidSizeBytes uint64 `json:"limit_void_size_bytes,omitempty"`
}

func (p *SConfigSettings) GetMessageSizeBytes() uint64 {
	return p.FMessageSizeBytes
}

func (p *SConfigSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *SConfigSettings) GetKeySizeBits() uint64 {
	return p.FKeySizeBits
}

func (p *SConfigSettings) GetQueuePeriodMS() uint64 {
	return p.FQueuePeriodMS
}

func (p *SConfigSettings) GetLimitVoidSizeBytes() uint64 {
	return p.FLimitVoidSizeBytes
}
