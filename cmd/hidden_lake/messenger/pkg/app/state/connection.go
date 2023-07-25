package state

var (
	_ IConnection = &SConnection{}
)

type SConnection struct {
	FAddress  string `json:"address"`
	FIsBackup bool   `json:"is_backup"`
}

func (p *SConnection) GetAddress() string {
	return p.FAddress
}

func (p *SConnection) IsBackup() bool {
	return p.FIsBackup
}
