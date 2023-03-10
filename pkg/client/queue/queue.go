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
	_ IMessageQueue = &sMessageQueue{}
)

type sMessageQueue struct {
	fIsRun    bool
	fMutex    sync.Mutex
	fSettings ISettings
	fClient   client.IClient
	fQueue    chan message.IMessage
	fMsgPull  sPull
}

type sPull struct {
	fSignal chan struct{}
	fQueue  chan message.IMessage
}

func NewMessageQueue(sett ISettings, client client.IClient) IMessageQueue {
	return &sMessageQueue{
		fSettings: sett,
		fClient:   client,
		fQueue:    make(chan message.IMessage, sett.GetCapacity()),
		fMsgPull: sPull{
			fQueue: make(chan message.IMessage, sett.GetPullCapacity()),
		},
	}
}

func (q *sMessageQueue) GetSettings() ISettings {
	return q.fSettings
}

func (q *sMessageQueue) UpdateClient(c client.IClient) {
	q.fMutex.Lock()
	defer q.fMutex.Unlock()

	q.fClient = c
	q.fQueue = make(chan message.IMessage, q.GetSettings().GetCapacity())
}

func (q *sMessageQueue) GetClient() client.IClient {
	q.fMutex.Lock()
	defer q.fMutex.Unlock()

	return q.fClient
}

func (q *sMessageQueue) Run() error {
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
			case <-q.readSignal():
				return
			case <-time.After(q.GetSettings().GetDuration() / 2):
				currLen := len(q.fMsgPull.fQueue)
				if uint64(currLen) >= q.GetSettings().GetPullCapacity() {
					continue
				}
				q.fMsgPull.fQueue <- q.newPseudoMessage()
			}
		}
	}()

	return nil
}

func (q *sMessageQueue) Stop() error {
	q.fMutex.Lock()
	defer q.fMutex.Unlock()

	if !q.fIsRun {
		return errors.New("queue already closed or not started")
	}
	q.fIsRun = false

	close(q.fMsgPull.fSignal)
	return nil
}

func (q *sMessageQueue) EnqueueMessage(msg message.IMessage) error {
	q.fMutex.Lock()
	defer q.fMutex.Unlock()

	if uint64(len(q.fQueue)) >= q.GetSettings().GetCapacity() {
		return errors.New("queue already full, need wait and retry")
	}

	q.fQueue <- msg
	return nil
}

func (q *sMessageQueue) DequeueMessage() <-chan message.IMessage {
	closed := make(chan bool)

	go func() {
		select {
		case <-q.readSignal():
			closed <- true
			return
		case <-time.After(q.GetSettings().GetDuration()):
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

func (q *sMessageQueue) newPseudoMessage() message.IMessage {
	msg, err := q.GetClient().EncryptPayload(
		q.GetClient().GetPubKey(),
		payload.NewPayload(0, []byte{1}),
	)
	if err != nil {
		panic(err)
	}
	return msg
}

func (q *sMessageQueue) readSignal() <-chan struct{} {
	q.fMutex.Lock()
	defer q.fMutex.Unlock()

	return q.fMsgPull.fSignal
}
