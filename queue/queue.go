package queue

import (
	"errors"
	"sync"
	"time"

	"github.com/number571/go-peer/client"
	"github.com/number571/go-peer/crypto/random"
	"github.com/number571/go-peer/message"
	"github.com/number571/go-peer/payload"
	"github.com/number571/go-peer/routing"
)

type sQueue struct {
	fMutex    sync.Mutex
	fSettings ISettings
	fClient   client.IClient
	fEnqueue  chan message.IMessage
	fMsgPull  *sPull
}

type sPull struct {
	fEnable  bool
	fSignal  chan struct{}
	fEnqueue chan message.IMessage
}

func NewQueue(sett ISettings, client client.IClient) IQueue {
	return &sQueue{
		fSettings: sett,
		fClient:   client,
		fEnqueue:  make(chan message.IMessage, sett.GetMainCapacity()),
		fMsgPull: &sPull{
			fSignal:  make(chan struct{}),
			fEnqueue: make(chan message.IMessage, sett.GetPullCapacity()),
		},
	}
}

func (q *sQueue) Settings() ISettings {
	return q.fSettings
}

func (q *sQueue) Start() error {
	if !q.runFullPull() {
		return errors.New("pull already enabled")
	}
	return nil
}

func (q *sQueue) Close() error {
	q.fMutex.Lock()
	defer q.fMutex.Unlock()

	if !q.fMsgPull.fEnable {
		return errors.New("pull already disabled")
	}

	q.fMsgPull.fSignal <- struct{}{}
	return nil
}

func (q *sQueue) Enqueue(n uint64, msg message.IMessage) error {
	if uint64(len(q.fEnqueue)) == q.Settings().GetMainCapacity() {
		return errors.New("queue already full, need wait and retry")
	}

	go func() {
		q.fMutex.Lock()
		defer q.fMutex.Unlock()

		// append pseudo messages to enqueue
		for i := uint64(0); i < n; i++ {
			q.fEnqueue <- (<-q.fMsgPull.fEnqueue)
		}

		// append main message to enqueue
		q.fEnqueue <- msg
	}()

	return nil
}

func (q *sQueue) Dequeue() <-chan message.IMessage {
	time.Sleep(q.Settings().GetDuration())

	go func() {
		if len(q.fEnqueue) == 0 {
			q.fEnqueue <- (<-q.fMsgPull.fEnqueue)
		}
	}()

	return q.fEnqueue
}

func (q *sQueue) runFullPull() bool {
	q.fMutex.Lock()
	defer q.fMutex.Unlock()

	if q.fMsgPull.fEnable {
		return false
	}
	q.fMsgPull.fEnable = true

	go func() {
		for {
			select {
			case <-q.fMsgPull.fSignal:
				q.fMsgPull.fEnable = false
				return
			default:
				currLen := len(q.fMsgPull.fEnqueue)
				if uint64(currLen) == q.Settings().GetPullCapacity() {
					time.Sleep(q.Settings().GetDuration())
					continue
				}
				q.fMsgPull.fEnqueue <- q.newPseudoMessage()
			}
		}
	}()

	return true
}

func (q *sQueue) newPseudoMessage() message.IMessage {
	rand := random.NewStdPRNG()
	return q.fClient.Encrypt(
		routing.NewRoute(q.fClient.PubKey()),
		payload.NewPayload(0, rand.Bytes(q.calcRandSize())),
	)
}

func (q *sQueue) calcRandSize() uint64 {
	msgSize := q.Settings().GetMessageSize()
	return random.NewStdPRNG().Uint64() % (msgSize / 2)
}
