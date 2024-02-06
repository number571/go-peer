package msgbroker

import "github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"

type IMessageBroker interface {
	Produce(string, utils.SMessage)
	Consume(string) (utils.SMessage, bool)
}
