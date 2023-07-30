package settings

var (
	_ IConfigSettings = &SConfigSettings{}
)

type SConfigSettings struct {
	FSettings SConfigSettingsBlock `json:"settings"`
}

type SConfigSettingsBlock struct {
	// basic values
	FMessageSizeBytes uint64 `json:"message_size_bytes"`
	FWorkSizeBits     uint64 `json:"work_size_bits"`

	// HLS, HLT
	FQueuePeriodMS      uint64 `json:"queue_period_ms,omitempty"`
	FLimitVoidSizeBytes uint64 `json:"limit_void_size_bytes,omitempty"`

	// HLS, HLM
	FKeySizeBits uint64 `json:"key_size_bits,omitempty"`

	// HLT, HLM
	FMessagesCapacity uint64 `json:"messages_capacity,omitempty"`
}

func (p *SConfigSettings) IsValid() bool {
	return p.FSettings.FMessageSizeBytes != 0 && p.FSettings.FWorkSizeBits != 0
}

func (p *SConfigSettings) GetMessageSizeBytes() uint64 {
	return p.FSettings.FMessageSizeBytes
}

func (p *SConfigSettings) GetWorkSizeBits() uint64 {
	return p.FSettings.FWorkSizeBits
}

func (p *SConfigSettings) GetKeySizeBits() uint64 {
	return p.FSettings.FKeySizeBits
}

func (p *SConfigSettings) GetQueuePeriodMS() uint64 {
	return p.FSettings.FQueuePeriodMS
}

func (p *SConfigSettings) GetMessagesCapacity() uint64 {
	return p.FSettings.FMessagesCapacity
}

func (p *SConfigSettings) GetLimitVoidSizeBytes() uint64 {
	return p.FSettings.FLimitVoidSizeBytes
}
