package queue

import (
	"errors"
	"sync"
	"time"

	"github.com/number571/go-peer/modules/client"
	"github.com/number571/go-peer/modules/crypto/random"
	"github.com/number571/go-peer/modules/message"
	"github.com/number571/go-peer/modules/payload"
)

type sQueue struct {
	fMutex    sync.Mutex
	fSettings ISettings
	fClient   client.IClient
	fQueue    chan message.IMessage
	fMsgPull  *sPull
}

type sPull struct {
	fEnable bool
	fSignal chan struct{}
	fQueue  chan message.IMessage
}

func NewQueue(sett ISettings, client client.IClient) IQueue {
	return &sQueue{
		fSettings: sett,
		fClient:   client,
		fQueue:    make(chan message.IMessage, sett.GetCapacity()),
		fMsgPull: &sPull{
			fSignal: make(chan struct{}),
			fQueue:  make(chan message.IMessage, sett.GetPullCapacity()),
		},
	}
}

func (q *sQueue) Settings() ISettings {
	return q.fSettings
}

func (q *sQueue) Client() client.IClient {
	return q.fClient
}

func (q *sQueue) Run() error {
	q.fMutex.Lock()
	defer q.fMutex.Unlock()

	if q.fMsgPull.fEnable {
		return errors.New("pull already enabled")
	}
	q.fMsgPull.fEnable = true

	go func() {
		for {
			select {
			case <-q.fMsgPull.fSignal:
				q.fMsgPull.fEnable = false
				return
			default:
				currLen := len(q.fMsgPull.fQueue)
				if uint64(currLen) == q.Settings().GetPullCapacity() {
					time.Sleep(q.Settings().GetDuration())
					continue
				}
				q.fMsgPull.fQueue <- q.newPseudoMessage()
			}
		}
	}()

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

func (q *sQueue) Enqueue(msg message.IMessage) error {
	if uint64(len(q.fQueue)) == q.Settings().GetCapacity() {
		return errors.New("queue already full, need wait and retry")
	}

	go func() {
		q.fMutex.Lock()
		defer q.fMutex.Unlock()

		q.fQueue <- msg
	}()

	return nil
}

func (q *sQueue) Dequeue() <-chan message.IMessage {
	time.Sleep(q.Settings().GetDuration())

	go func() {
		q.fMutex.Lock()
		defer q.fMutex.Unlock()

		if len(q.fQueue) == 0 {
			q.fQueue <- (<-q.fMsgPull.fQueue)
		}
	}()

	return q.fQueue
}

func (q *sQueue) newPseudoMessage() message.IMessage {
	rand := random.NewStdPRNG()
	msg, err := q.fClient.Encrypt(
		q.fClient.PubKey(),
		payload.NewPayload(0, rand.Bytes(1)),
	)
	if err != nil {
		panic(err)
	}
	return msg
}
