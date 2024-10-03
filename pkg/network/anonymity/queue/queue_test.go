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

func TestPanicNewQBProblemProcessor(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	client := client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: testutils.TCMessageSize,
			FKeySizeBits:      testutils.TcKeySize,
		}),
		asymmetric.LoadRSAPrivKey(testutils.TcPrivKey1024),
	)
	_ = NewQBProblemProcessor(
		NewSettings(&SSettings{
			FMessageConstructSettings: net_message.NewConstructSettings(&net_message.SConstructSettings{
				FSettings: net_message.NewSettings(&net_message.SSettings{}),
			}),
			FMainPoolCapacity: testutils.TCQueueCapacity,
			FRandPoolCapacity: 1,
			FQueuePeriod:      100 * time.Millisecond,
		}),
		client,
		asymmetric.LoadRSAPrivKey(testutils.TcPrivKey2048).GetPubKey(),
	)
}

func TestSettings(t *testing.T) {
	t.Parallel()

	for i := 0; i < 3; i++ {
		testSettings(t, i)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	// check FMessageConstructSettings
	_ = (&sSettings{
		FMainPoolCapacity: testutils.TCQueueCapacity,
		FRandPoolCapacity: testutils.TCQueueCapacity,
		FQueuePeriod:      500 * time.Millisecond,
	}).mustNotNull()
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
	case 2:
		_ = NewSettings(&SSettings{
			FMainPoolCapacity: testutils.TCQueueCapacity,
			FRandPoolCapacity: testutils.TCQueueCapacity,
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
		asymmetric.LoadRSAPrivKey(testutils.TcPrivKey1024),
	)
	queue := NewQBProblemProcessor(
		NewSettings(&SSettings{
			FMessageConstructSettings: net_message.NewConstructSettings(&net_message.SConstructSettings{
				FSettings: net_message.NewSettings(&net_message.SSettings{}),
			}),
			FMainPoolCapacity: testutils.TCQueueCapacity,
			FRandPoolCapacity: 1,
			FQueuePeriod:      100 * time.Millisecond,
		}),
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
		sQueue := queue.(*sQBProblemProcessor)
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

	queue := NewQBProblemProcessor(
		NewSettings(&SSettings{
			FMessageConstructSettings: net_message.NewConstructSettings(&net_message.SConstructSettings{
				FSettings: net_message.NewSettings(&net_message.SSettings{
					FWorkSizeBits: 10,
				}),
			}),
			FNetworkMask:      1,
			FMainPoolCapacity: testutils.TCQueueCapacity,
			FRandPoolCapacity: testutils.TCQueueCapacity,
			FQueuePeriod:      100 * time.Millisecond,
			FRandQueuePeriod:  100 * time.Millisecond,
		}),
		client.NewClient(
			message.NewSettings(&message.SSettings{
				FMessageSizeBytes: testutils.TCMessageSize,
				FKeySizeBits:      testutils.TcKeySize,
			}),
			asymmetric.LoadRSAPrivKey(testutils.TcPrivKey1024),
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

func testQueue(queue IQBProblemProcessor) error {
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

	// auto fill queue enabled only if QB=true
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
