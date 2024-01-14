package chat_queue

import (
	"testing"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"
)

func TestChatQueue(t *testing.T) {
	t.Parallel()

	chatQueue := NewChatQueue(3)

	chatQueue.Init()

	addr := "address"
	msgInfo := "msg_info"

	go func() {
		chatQueue.Push(&SMessage{
			FAddress: addr,
			FMessageInfo: utils.SMessageInfo{
				FMessage: msgInfo,
			},
		})
	}()

	msg, ok := chatQueue.Load(addr)
	if !ok {
		t.Error("got not ok load")
		return
	}

	if msg.FAddress != addr {
		t.Error("msg.FAddress != addr")
		return
	}

	if msg.FMessageInfo.FMessage != msgInfo {
		t.Error("msg.FMessageInfo.FMessage != msgInfo")
		return
	}
}
