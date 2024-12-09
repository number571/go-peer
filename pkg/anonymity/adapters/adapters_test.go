package adapters

import (
	"bytes"
	"context"
	"testing"

	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
)

const (
	tcMessage = "hello, world!"
)

func TestAdapter(t *testing.T) {
	msgChan := make(chan message.IMessage, 1)
	adapter := NewAdapterByFuncs(
		func(_ context.Context, msg message.IMessage) error {
			msgChan <- msg
			return nil
		},
		func(_ context.Context) (message.IMessage, error) {
			return <-msgChan, nil
		},
	)

	ctx := context.Background()

	err := adapter.Produce(ctx, message.NewMessage(
		message.NewConstructSettings(&message.SConstructSettings{
			FSettings: message.NewSettings(&message.SSettings{}),
		}),
		payload.NewPayload32(0x01, []byte(tcMessage)),
	))
	if err != nil {
		t.Error(err)
		return
	}

	msg, err := adapter.Consume(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(msg.GetPayload().GetBody(), []byte(tcMessage)) {
		t.Error("consume invalid message")
		return
	}
}
