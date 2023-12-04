package queue

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/_data"
)

func TestSettings(t *testing.T) {
	t.Parallel()

	for i := 0; i < 3; i++ {
		testSettings(t, i)
	}
}

func testSettings(t *testing.T, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	switch n {
	case 0:
		_ = NewSettings(&SSettings{
			FPoolCapacity: testutils.TCQueueCapacity,
			FDuration:     500 * time.Millisecond,
		})
	case 1:
		_ = NewSettings(&SSettings{
			FMainCapacity: testutils.TCQueueCapacity,
			FDuration:     500 * time.Millisecond,
		})
	case 2:
		_ = NewSettings(&SSettings{
			FMainCapacity: testutils.TCQueueCapacity,
			FPoolCapacity: testutils.TCQueueCapacity,
		})
	}
}

func TestRunStopQueue(t *testing.T) {
	t.Parallel()

	client := client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FKeySizeBits:      testutils.TcKeySize,
		}),
		asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024),
	)
	queue := NewMessageQueue(
		NewSettings(&SSettings{
			FMainCapacity: testutils.TCQueueCapacity,
			FPoolCapacity: 1,
			FDuration:     100 * time.Millisecond,
		}),
		client,
	)

	if err := queue.Run(); err != nil {
		t.Error(err)
		return
	}
	if err := queue.Run(); err == nil {
		t.Error("success run already running queue")
		return
	}

	time.Sleep(100 * time.Millisecond)

	msg, err := client.EncryptPayload(
		client.GetPubKey(),
		payload.NewPayload(0, []byte(testutils.TcBody)),
	)
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < testutils.TCQueueCapacity; i++ {
		if err := queue.EnqueueMessage(msg); err != nil {
			t.Error(err)
			return
		}
	}

	if err := queue.EnqueueMessage(msg); err == nil {
		t.Error("success enqueue message with max capacity")
		return
	}

	if err := queue.Stop(); err != nil {
		t.Error(err)
		return
	}
	if err := queue.Stop(); err == nil {
		t.Error("success stop already stopped queue")
		return
	}
}

func TestQueue(t *testing.T) {
	t.Parallel()

	oldClient := client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FKeySizeBits:      testutils.TcKeySize,
		}),
		asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024),
	)
	queue := NewMessageQueue(
		NewSettings(&SSettings{
			FMainCapacity: testutils.TCQueueCapacity,
			FPoolCapacity: testutils.TCQueueCapacity,
			FDuration:     100 * time.Millisecond,
		}),
		oldClient,
	).WithNetworkSettings(
		uint64(1),
		net_message.NewSettings(&net_message.SSettings{
			FNetworkKey:   "old_network_key",
			FWorkSizeBits: 10,
		}),
	)

	sett := queue.GetSettings()
	if sett.GetMainCapacity() != testutils.TCQueueCapacity {
		t.Error("sett.GetMainCapacity() != testutils.TCQueueCapacity")
		return
	}

	if err := testQueue(queue); err != nil {
		t.Error(err)
		return
	}
}

func testQueue(queue IMessageQueue) error {
	if err := queue.Run(); err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)

	queue.WithNetworkSettings(
		uint64(1),
		net_message.NewSettings(&net_message.SSettings{
			FNetworkKey:   "new_network_key",
			FWorkSizeBits: 10,
		}),
	)

	time.Sleep(100 * time.Millisecond)

	msgs := make([]net_message.IMessage, 0, 3)
	for i := 0; i < 3; i++ {
		msgs = append(msgs, <-queue.DequeueMessage())
	}

	for i := 0; i < len(msgs)-1; i++ {
		for j := i + 1; j < len(msgs); j++ {
			if bytes.Equal(msgs[i].GetHash(), msgs[j].GetHash()) {
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

	hash := msg.GetHash()
	for i := 0; i < 3; i++ {
		queue.EnqueueMessage(msg)
	}
	for i := 0; i < 3; i++ {
		netMsg := <-queue.DequeueMessage()
		msg, err := message.LoadMessage(client.GetSettings(), netMsg.GetPayload().GetBody())
		if err != nil {
			return err
		}
		if !bytes.Equal(msg.GetHash(), hash) {
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
