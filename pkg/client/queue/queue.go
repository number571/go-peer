package queue

import (
	"errors"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/payload"
)

var (
	_ IQueue = &sQueue{}
)

type sQueue struct {
	fIsRun    bool
	fMutex    sync.Mutex
	fSettings ISettings
	fClient   client.IClient
	fQueue    chan message.IMessage
	fMsgPull  *sPull
}

type sPull struct {
	fSignal chan struct{}
	fQueue  chan message.IMessage
}

func NewQueue(sett ISettings, client client.IClient) IQueue {
	return &sQueue{
		fSettings: sett,
		fClient:   client,
		fQueue:    make(chan message.IMessage, sett.GetCapacity()),
		fMsgPull: &sPull{
			fQueue: make(chan message.IMessage, sett.GetPullCapacity()),
		},
	}
}

func (q *sQueue) Settings() ISettings {
	return q.fSettings
}

func (q *sQueue) UpdateClient(c client.IClient) {
	q.fMutex.Lock()
	defer q.fMutex.Unlock()

	q.fClient = c
	q.fQueue = make(chan message.IMessage, q.Settings().GetCapacity())
}

func (q *sQueue) Client() client.IClient {
	return q.fClient
}

func (q *sQueue) Run() error {
	q.fMutex.Lock()
	defer q.fMutex.Unlock()

	if q.fIsRun {
		return errors.New("queue already running")
	}
	q.fIsRun = true

	q.fMsgPull.fSignal = make(chan struct{})
	go func() {
		for {
			select {
			case <-q.fMsgPull.fSignal:
				return
			case <-time.After(q.Settings().GetDuration() / 2):
				currLen := len(q.fMsgPull.fQueue)
				if uint64(currLen) >= q.Settings().GetPullCapacity() {
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

	if !q.fIsRun {
		return errors.New("queue already closed or not started")
	}
	q.fIsRun = false

	close(q.fMsgPull.fSignal)
	return nil
}

func (q *sQueue) Enqueue(msg message.IMessage) error {
	q.fMutex.Lock()
	defer q.fMutex.Unlock()

	if uint64(len(q.fQueue)) >= q.Settings().GetCapacity() {
		return errors.New("queue already full, need wait and retry")
	}

	q.fQueue <- msg
	return nil
}

func (q *sQueue) Dequeue() <-chan message.IMessage {
	closed := make(chan bool)

	go func() {
		select {
		case <-q.fMsgPull.fSignal:
			closed <- true
			return
		case <-time.After(q.Settings().GetDuration()):
			q.fMutex.Lock()
			defer q.fMutex.Unlock()

			if len(q.fQueue) == 0 {
				q.fQueue <- (<-q.fMsgPull.fQueue)
			}
			closed <- false
		}
	}()

	if <-closed {
		queue := make(chan message.IMessage)
		close(queue)
		return queue
	}

	return q.fQueue
}

func (q *sQueue) newPseudoMessage() message.IMessage {
	msg, err := q.Client().Encrypt(
		q.Client().PubKey(),
		payload.NewPayload(0, []byte{1}),
	)
	if err != nil {
		panic(err)
	}
	return msg
}
