package queue

import (
	"errors"
	"sync"
	"time"

	"github.com/number571/go-peer/modules/client"
	"github.com/number571/go-peer/modules/message"
	"github.com/number571/go-peer/modules/payload"
)

var (
	_ IQueue = &sQueue{}
)

type iState int

const (
	cIsInit iState = iota
	cIsRun
	cIsClose
)

type sQueue struct {
	fMutex    sync.Mutex
	fSettings ISettings
	fClient   client.IClient
	fQueue    chan message.IMessage
	fMsgPull  *sPull
}

type sPull struct {
	fState  iState
	fSignal chan struct{}
	fQueue  chan message.IMessage
}

func NewQueue(sett ISettings, client client.IClient) IQueue {
	return &sQueue{
		fSettings: sett,
		fClient:   client,
		fQueue:    make(chan message.IMessage, sett.GetCapacity()),
		fMsgPull: &sPull{
			fState:  cIsInit,
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

	if q.fMsgPull.fState != cIsInit {
		return errors.New("queue already started or closed")
	}
	q.fMsgPull.fState = cIsRun

	go func() {
		for {
			select {
			case <-q.fMsgPull.fSignal:
				q.fMsgPull.fState = cIsClose
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

	if q.fMsgPull.fState != cIsRun {
		return errors.New("queue already closed or not started")
	}
	q.fMsgPull.fSignal <- struct{}{}

	close(q.fQueue)
	close(q.fMsgPull.fQueue)

	return nil
}

func (q *sQueue) Enqueue(msg message.IMessage) error {
	q.fMutex.Lock()
	defer q.fMutex.Unlock()

	if uint64(len(q.fQueue)) == q.Settings().GetCapacity() {
		return errors.New("queue already full, need wait and retry")
	}

	go func() {
		q.fMutex.Lock()
		defer q.fMutex.Unlock()

		if q.fMsgPull.fState != cIsRun {
			return
		}

		q.fQueue <- msg
	}()

	return nil
}

func (q *sQueue) Dequeue() <-chan message.IMessage {
	time.Sleep(q.Settings().GetDuration())

	go func() {
		q.fMutex.Lock()
		defer q.fMutex.Unlock()

		if q.fMsgPull.fState != cIsRun {
			return
		}

		if len(q.fQueue) == 0 {
			q.fQueue <- (<-q.fMsgPull.fQueue)
		}
	}()

	return q.fQueue
}

func (q *sQueue) newPseudoMessage() message.IMessage {
	msg, err := q.fClient.Encrypt(
		q.fClient.PubKey(),
		payload.NewPayload(0, []byte{1}),
	)
	if err != nil {
		panic(err)
	}
	return msg
}
