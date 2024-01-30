package receiver

import "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"

type SMessage struct {
	FAddress     string             `json:"address"`
	FMessageInfo utils.SMessageInfo `json:"message_info"`
}

type IMessageReceiver interface {
	Init(string) IMessageReceiver
	Send(*SMessage)
	Recv() (*SMessage, bool)
}
