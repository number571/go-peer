package chat_queue

import "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"

type SMessage struct {
	FAddress     string             `json:"address"`
	FMessageInfo utils.SMessageInfo `json:"message_info"`
}

type IChatQueue interface {
	Init()
	Push(*SMessage)
	Load(string) (*SMessage, bool)
}
