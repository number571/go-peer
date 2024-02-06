package msgbroker

import (
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"
)

func TestMessageBroker(t *testing.T) {
	t.Parallel()

	addr := "address"
	msgData := "msg_data"

	msgReceiver := NewMessageBroker()

	go func() {
		time.Sleep(100 * time.Millisecond)
		msgReceiver.Produce(addr, utils.SMessage{FMainData: msgData})
	}()

	msg, ok := msgReceiver.Consume(addr)
	if !ok {
		t.Error("got not ok recv")
		return
	}

	if msg.FMainData != msgData {
		t.Error("msg.FMainData != msgData")
		return
	}
}
