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

	for i := 0; i < 2; i++ {
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
			FRandPoolCapacity: testutils.TCQueueCapacity,
			FQueuePeriod:      500 * time.Millisecond,
		})
	case 1:
		_ = NewSettings(&SSettings{
			FMainPoolCapacity: testutils.TCQueueCapacity,
			FQueuePeriod:      500 * time.Millisecond,
		})
	}
}

func TestQueueVoidDisabled(t *testing.T) {
	t.Parallel()

	queue := NewMessageQueueProcessor(
		NewSettings(&SSettings{
			FNetworkMask:      1,
			FWorkSizeBits:     10,
			FMainPoolCapacity: testutils.TCQueueCapacity,
			FRandPoolCapacity: testutils.TCQueueCapacity,
			FParallel:         1,
			FRandQueuePeriod:  100 * time.Millisecond,
		}),
		NewVSettings(&SVSettings{
			FNetworkKey: "network_key",
		}),
		client.NewClient(
			message.NewSettings(&message.SSettings{
				FMessageSizeBytes: testutils.TCMessageSize,
				FKeySizeBits:      testutils.TcKeySize,
			}),
			asymmetric.LoadRSAPrivKey(testutils.Tc1PrivKey1024),
		),
		asymmetric.NewRSAPrivKey(testutils.TcKeySize).GetPubKey(),
	)

	if err := testQueue(queue); err != nil {
		t.Error(err)
		return
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
	queue := NewMessageQueueProcessor(
		NewSettings(&SSettings{
			FMainPoolCapacity: testutils.TCQueueCapacity,
			FRandPoolCapacity: 1,
			FParallel:         1,
			FQueuePeriod:      100 * time.Millisecond,
		}),
		NewVSettings(&SVSettings{}),
		client,
		asymmetric.NewRSAPrivKey(client.GetPrivKey().GetSize()).GetPubKey(),
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
		sQueue := queue.(*sMessageQueueProcessor)
		if len(sQueue.fRandPool.fQueue) == int(sett.GetRandPoolCapacity()) {
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

	pubKey := client.GetPubKey()
	pldBytes := payload.NewPayload64(0, []byte(testutils.TcBody)).ToBytes()
	for i := 0; i < testutils.TCQueueCapacity; i++ {
		if err := queue.EnqueueMessage(pubKey, pldBytes); err != nil {
			t.Error(err)
			return
		}
	}

	// after full queue
	for i := 0; i < 2*testutils.TCQueueCapacity; i++ {
		if err := queue.EnqueueMessage(pubKey, pldBytes); err != nil {
			return
		}
	}

	t.Error("success enqueue message with max capacity")
}

func TestQueue(t *testing.T) {
	t.Parallel()

	queue := NewMessageQueueProcessor(
		NewSettings(&SSettings{
			FNetworkMask:      1,
			FWorkSizeBits:     10,
			FMainPoolCapacity: testutils.TCQueueCapacity,
			FRandPoolCapacity: testutils.TCQueueCapacity,
			FParallel:         1,
			FQueuePeriod:      100 * time.Millisecond,
			FRandQueuePeriod:  100 * time.Millisecond,
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
		asymmetric.NewRSAPrivKey(testutils.TcKeySize).GetPubKey(),
	)

	sett := queue.GetSettings()
	if sett.GetMainPoolCapacity() != testutils.TCQueueCapacity {
		t.Error("sett.GetMainCapacity() != testutils.TCQueueCapacity")
		return
	}

	if err := testQueue(queue); err != nil {
		t.Error(err)
		return
	}
}

func testQueue(queue IMessageQueueProcessor) error {
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
	pubKey := client.GetPubKey()
	pldBytes := payload.NewPayload64(0, []byte(testutils.TcBody)).ToBytes()
	if err := queue.EnqueueMessage(pubKey, pldBytes); err != nil {
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

	// auto fill queue enabled only if QB=true
	if queue.GetSettings().GetQueuePeriod() != 0 {
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
