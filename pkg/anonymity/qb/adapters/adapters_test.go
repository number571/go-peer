package adapters

import (
	"bytes"
	"context"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
)

const (
	tcMessage = "hello, world!"
)

func TestAdapter(t *testing.T) {
	msgChan := make(chan layer1.IMessage, 1)
	adapter := NewAdapterByFuncs(
		func(_ context.Context, msg layer1.IMessage) error {
			msgChan <- msg
			return nil
		},
		func(_ context.Context) (layer1.IMessage, error) {
			return <-msgChan, nil
		},
	)

	ctx := context.Background()

	err := adapter.Produce(ctx, layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{}),
		}),
		[]byte(tcMessage),
	))
	if err != nil {
		t.Fatal(err)
	}

	msg, err := adapter.Consume(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(msg.GetBody(), []byte(tcMessage)) {
		t.Fatal("consume invalid message")
	}
}
