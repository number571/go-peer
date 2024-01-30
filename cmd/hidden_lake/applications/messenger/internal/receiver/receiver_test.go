package receiver

import (
	"testing"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/messenger/internal/utils"
)

func TestMessageReceiver(t *testing.T) {
	t.Parallel()

	addr := "address"
	msgInfo := "msg_info"

	msgReceiver := NewMessageReceiver().Init(addr)

	go func() {
		msgReceiver.Send(&SMessage{
			FAddress: addr,
			FMessageInfo: utils.SMessageInfo{
				FMessage: msgInfo,
			},
		})
	}()

	msg, ok := msgReceiver.Recv()
	if !ok {
		t.Error("got not ok recv")
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
