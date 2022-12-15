package queue

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/number571/go-peer/internal/testutils"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/payload"
)

func TestQueue(t *testing.T) {
	oldClient := client.NewClient(
		client.NewSettings(&client.SSettings{
			FWorkSize:    10,
			FMessageSize: (1 << 20),
		}),
		asymmetric.LoadRSAPrivKey(testutils.TcPrivKey),
	)
	queue := NewQueue(
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

func testQueue(queue IQueue) error {
	if err := queue.Run(); err != nil {
		return err
	}

	msgs := make([]message.IMessage, 0, 3)
	for i := 0; i < 3; i++ {
		msgs = append(msgs, <-queue.Dequeue())
	}

	for i := 0; i < len(msgs)-1; i++ {
		for j := i + 1; j < len(msgs); j++ {
			if bytes.Equal(msgs[i].Body().Hash(), msgs[j].Body().Hash()) {
				return fmt.Errorf("hash of messages equals (%d and %d)", i, i)
			}
		}
	}

	msg, err := queue.Client().Encrypt(
		queue.Client().PubKey(),
		payload.NewPayload(0, []byte(testutils.TcBody)),
	)
	if err != nil {
		return err
	}

	hash := msg.Body().Hash()
	for i := 0; i < 3; i++ {
		queue.Enqueue(msg)
	}
	for i := 0; i < 3; i++ {
		msg := <-queue.Dequeue()
		if !bytes.Equal(msg.Body().Hash(), hash) {
			return fmt.Errorf("hash of messages not equals (%d)", i)
		}
	}

	if err := queue.Close(); err != nil {
		return err
	}

	return nil
}
