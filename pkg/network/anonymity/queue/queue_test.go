// nolint: goerr113
package queue

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SQueueError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

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
			FVoidCapacity: testutils.TCQueueCapacity,
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
			FVoidCapacity: testutils.TCQueueCapacity,
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
			FVoidCapacity: 1,
			FParallel:     1,
			FDuration:     100 * time.Millisecond,
		}),
		NewVSettings(&SVSettings{}),
		client,
	)

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	go func() {
		if err := queue.Run(ctx1); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
			return
		}
	}()

	err := testutils.TryN(50, 10*time.Millisecond, func() error {
		sett := queue.GetSettings()
		sQueue := queue.(*sMessageQueue)
		if len(sQueue.fVoidPool.fQueue) == int(sett.GetVoidCapacity()) {
			return nil
		}
		return errors.New("len(void queue) != max capacity")
	})
	if err != nil {
		t.Error(err)
		return
	}

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	go func() {
		if err := queue.Run(ctx2); err == nil {
			t.Error("success run already running queue")
			return
		}
	}()

	msg, err := client.EncryptPayload(
		client.GetPubKey(),
		payload.NewPayload64(0, []byte(testutils.TcBody)),
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

	// after full queue
	for i := 0; i < 2*testutils.TCQueueCapacity; i++ {
		if err := queue.EnqueueMessage(msg); err != nil {
			return
		}
	}

	t.Error("success enqueue message with max capacity")
}

func TestQueue(t *testing.T) {
	t.Parallel()

	queue := NewMessageQueue(
		NewSettings(&SSettings{
			FNetworkMask:  1,
			FWorkSizeBits: 10,
			FMainCapacity: testutils.TCQueueCapacity,
			FVoidCapacity: testutils.TCQueueCapacity,
			FParallel:     1,
			FDuration:     100 * time.Millisecond,
			FRandDuration: 100 * time.Millisecond,
		}),
		NewVSettings(&SVSettings{
			FNetworkKey: "old_network_key",
		}),
		client.NewClient(
			message.NewSettings(&message.SSettings{
				FMessageSizeBytes: testutils.TCMessageSize,
				FKeySizeBits:      testutils.TcKeySize,
			}),
			asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024),
		),
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
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		time.Sleep(200 * time.Millisecond)
	}()

	go func() {
		if err := queue.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			return
		}
	}()

	client := queue.GetClient()
	msg, err := client.EncryptPayload(
		client.GetPubKey(),
		payload.NewPayload64(0, []byte(testutils.TcBody)),
	)
	if err != nil {
		return err
	}

	if err := queue.EnqueueMessage(msg); err != nil {
		return err
	}

	// wait minimum one generated message
	time.Sleep(300 * time.Millisecond)

	// clear old messages
	newNetworkKey := "new_network_key"
	queue.SetVSettings(NewVSettings(&SVSettings{
		FNetworkKey: newNetworkKey,
	}))

	nVSettings := queue.GetVSettings()
	if nVSettings.GetNetworkKey() != newNetworkKey {
		return errors.New("incorrect set variable settings")
	}

	msgs := make([]net_message.IMessage, 0, 3)
	for i := 0; i < 3; i++ {
		msgs = append(msgs, queue.DequeueMessage(ctx))
	}

	for i := 0; i < len(msgs)-1; i++ {
		for j := i + 1; j < len(msgs); j++ {
			if bytes.Equal(msgs[i].GetHash(), msgs[j].GetHash()) {
				return fmt.Errorf("hash of messages equals (%d and %d)", i, i)
			}
		}
	}

	msg2, err := client.EncryptPayload(
		client.GetPubKey(),
		payload.NewPayload64(0, []byte(testutils.TcBody)),
	)
	if err != nil {
		return err
	}

	hash := msg2.GetEnck()
	for i := 0; i < 3; i++ {
		if err := queue.EnqueueMessage(msg2); err != nil {
			return err
		}
	}

	for i := 0; i < 3; i++ {
		netMsg := queue.DequeueMessage(ctx)
		msg, err := message.LoadMessage(client.GetSettings(), netMsg.GetPayload().GetBody())
		if err != nil {
			return err
		}
		if !bytes.Equal(msg.GetEnck(), hash) {
			return fmt.Errorf("enc_key of messages not equals (%d)", i)
		}
	}

	notClosed := make(chan bool)
	go func() {
		// test close with parallel dequeue
		msg := queue.DequeueMessage(ctx)
		notClosed <- (msg != nil)
	}()

	cancel()
	if <-notClosed {
		return errors.New("success dequeue with close")
	}
	return nil
}
