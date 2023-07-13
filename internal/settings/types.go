package settings

import "github.com/number571/go-peer/pkg/client/message"

var (
	_ IConfigSettings = &SConfigSettings{}
)

type IConfigSettings interface {
	message.ISettings
}

type SConfigSettings struct {
	FSettings SConfigSettingsBlock `json:"settings"`
}

type SConfigSettingsBlock struct {
	FMessageSize uint64 `json:"message_size"`
	FWorkSize    uint64 `json:"work_size"`
}

func (p *SConfigSettings) GetMessageSize() uint64 {
	return p.FSettings.FMessageSize
}

func (p *SConfigSettings) GetWorkSize() uint64 {
	return p.FSettings.FWorkSize
}

func (p *SConfigSettingsBlock) IsValid() bool {
	return p.FMessageSize != 0 && p.FWorkSize != 0
}
