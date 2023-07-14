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

	// HLS, HLM
	FKeySizeBits uint64 `json:"key_size_bits,omitempty"`

	// HLS
	FQueuePeriodMS uint64 `json:"queue_period_ms,omitempty"`

	// HLT
	FMessagesCapacity uint64 `json:"messages_capacity,omitempty"`
}

func (p *SConfigSettingsBlock) IsValid() bool {
	return p.FMessageSizeBytes != 0 && p.FWorkSizeBits != 0
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
