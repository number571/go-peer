package chat_queue

type SMessage struct {
	FAddress   string `json:"address"`
	FMessage   string `json:"message"`
	FTimestamp string `json:"timestamp"`
}

type IChatQueue interface {
	Init()
	Push(*SMessage)
	Load(string) (*SMessage, bool)
}
