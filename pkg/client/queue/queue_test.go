package queue

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestQueue(t *testing.T) {
	oldClient := client.NewClient(
		client.NewSettings(&client.SSettings{
			FWorkSize:    10,
			FMessageSize: (1 << 20),
		}),
		asymmetric.LoadRSAPrivKey(testutils.TcPrivKey),
	)
	queue := NewMessageQueue(
		NewSettings(&SSettings{
			FCapacity:     10,
			FPullCapacity: 5,
			FDuration:     500 * time.Millisecond,
		}),
		oldClient,
	)

	if err := testQueue(queue); err != nil {
		t.Error(err)
		return
	}

	newClient := client.NewClient(
		client.NewSettings(&client.SSettings{
			FWorkSize:    10,
			FMessageSize: (1 << 20),
		}),
		asymmetric.LoadRSAPrivKey(testutils.TcPrivKey2),
	)
	queue.UpdateClient(newClient)

	if err := testQueue(queue); err != nil {
		t.Error(err)
		return
	}
}

func testQueue(queue IMessageQueue) error {
	if err := queue.Run(); err != nil {
		return err
	}

	msgs := make([]message.IMessage, 0, 3)
	for i := 0; i < 3; i++ {
		msgs = append(msgs, <-queue.DequeueMessage())
	}

	for i := 0; i < len(msgs)-1; i++ {
		for j := i + 1; j < len(msgs); j++ {
			if bytes.Equal(msgs[i].GetBody().GetHash(), msgs[j].GetBody().GetHash()) {
				return fmt.Errorf("hash of messages equals (%d and %d)", i, i)
			}
		}
	}

	client := queue.GetClient()
	msg, err := client.EncryptPayload(
		client.GetPubKey(),
		payload.NewPayload(0, []byte(testutils.TcBody)),
	)
	if err != nil {
		return err
	}

	hash := msg.GetBody().GetHash()
	for i := 0; i < 3; i++ {
		queue.EnqueueMessage(msg)
	}
	for i := 0; i < 3; i++ {
		msg := <-queue.DequeueMessage()
		if !bytes.Equal(msg.GetBody().GetHash(), hash) {
			return fmt.Errorf("hash of messages not equals (%d)", i)
		}
	}

	closed := make(chan bool)
	go func() {
		// test close with parallel dequeue
		_, ok := <-queue.DequeueMessage()
		closed <- ok
	}()

	if err := queue.Stop(); err != nil {
		return err
	}

	if <-closed {
		return fmt.Errorf("success dequeue with close")
	}
	return nil
}
