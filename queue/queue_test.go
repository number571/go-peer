package queue

import (
	"bytes"
	"testing"
	"time"

	"github.com/number571/go-peer/client"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/message"
	"github.com/number571/go-peer/payload"
	"github.com/number571/go-peer/routing"
	"github.com/number571/go-peer/testutils"
)

func TestQueue(t *testing.T) {
	privKey := asymmetric.LoadRSAPrivKey(testutils.TcPrivKey)
	client := client.NewClient(client.NewSettings(10, (1<<10)), privKey)
	queue := NewQueue(NewSettings(10, 5, (1<<20), 500*time.Millisecond), client)

	if err := queue.Start(); err != nil {
		t.Error(err)
		return
	}

	msgs := make([]message.IMessage, 0, 3)
	for i := 0; i < 3; i++ {
		msgs = append(msgs, <-queue.Dequeue())
	}

	for i := 0; i < len(msgs)-1; i++ {
		for j := i + 1; j < len(msgs); j++ {
			if bytes.Equal(msgs[i].Body().Hash(), msgs[j].Body().Hash()) {
				t.Errorf("hash of messages equals (%d and %d)", i, i)
				return
			}
		}
	}

	msg := client.Encrypt(
		routing.NewRoute(client.PubKey()),
		payload.NewPayload(0, []byte(testutils.TcBody)),
	)
	hash := msg.Body().Hash()

	for i := 0; i < 3; i++ {
		queue.Enqueue(0, msg)
	}
	for i := 0; i < 3; i++ {
		msg := <-queue.Dequeue()
		if !bytes.Equal(msg.Body().Hash(), hash) {
			t.Errorf("hash of messages not equals (%d)", i)
			return
		}
	}

	if err := queue.Close(); err != nil {
		t.Error(err)
		return
	}
}