package settings

var (
	_ IConfigSettings = &SConfigSettings{}
)

type SConfigSettings struct {
	FSettings SConfigSettingsBlock `json:"settings"`
}

type SConfigSettingsBlock struct {
	// basic values
	FMessageSize uint64 `json:"message_size"`
	FWorkSize    uint64 `json:"work_size"`

	// HLS, HLM
	FKeySize uint64 `json:"key_size,omitempty"`

	// HLS
	FQueuePeriod uint64 `json:"queue_period,omitempty"`

	// HLT
	FCapMessages uint64 `json:"cap_messages,omitempty"`
}

func (p *SConfigSettingsBlock) IsValid() bool {
	return p.FMessageSize != 0 && p.FWorkSize != 0
}

func (p *SConfigSettings) GetMessageSize() uint64 {
	return p.FSettings.FMessageSize
}

func (p *SConfigSettings) GetWorkSize() uint64 {
	return p.FSettings.FWorkSize
}

func (p *SConfigSettings) GetKeySize() uint64 {
	return p.FSettings.FKeySize
}

func (p *SConfigSettings) GetQueuePeriod() uint64 {
	return p.FSettings.FQueuePeriod
}

func (p *SConfigSettings) GetCapMessages() uint64 {
	return p.FSettings.FCapMessages
}
