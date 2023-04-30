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

func NewMessageQueue(pSett ISettings, pClient client.IClient) IMessageQueue {
	return &sMessageQueue{
		fSettings: pSett,
		fClient:   pClient,
		fQueue:    make(chan message.IMessage, pSett.GetCapacity()),
		fMsgPull: sPull{
			fQueue: make(chan message.IMessage, pSett.GetPullCapacity()),
		},
	}
}

func (p *sMessageQueue) GetSettings() ISettings {
	return p.fSettings
}

func (p *sMessageQueue) UpdateClient(pClient client.IClient) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fClient = pClient
	p.fQueue = make(chan message.IMessage, p.GetSettings().GetCapacity())
}

func (p *sMessageQueue) GetClient() client.IClient {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fClient
}

func (p *sMessageQueue) Run() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fIsRun {
		return errors.New("queue already running")
	}
	p.fIsRun = true

	p.fMsgPull.fSignal = make(chan struct{})
	go func() {
		for {
			select {
			case <-p.readSignal():
				return
			case <-time.After(p.GetSettings().GetDuration() / 2):
				currLen := len(p.fMsgPull.fQueue)
				if uint64(currLen) >= p.GetSettings().GetPullCapacity() {
					continue
				}
				p.fMsgPull.fQueue <- p.newPseudoMessage()
			}
		}
	}()

	return nil
}

func (p *sMessageQueue) Stop() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if !p.fIsRun {
		return errors.New("queue already closed or not started")
	}
	p.fIsRun = false

	close(p.fMsgPull.fSignal)
	return nil
}

func (p *sMessageQueue) EnqueueMessage(pMsg message.IMessage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if uint64(len(p.fQueue)) >= p.GetSettings().GetCapacity() {
		return errors.New("queue already full, need wait and retry")
	}

	p.fQueue <- pMsg
	return nil
}

func (p *sMessageQueue) DequeueMessage() <-chan message.IMessage {
	closed := make(chan bool)

	go func() {
		select {
		case <-p.readSignal():
			closed <- true
			return
		case <-time.After(p.GetSettings().GetDuration()):
			p.fMutex.Lock()
			defer p.fMutex.Unlock()

			if len(p.fQueue) == 0 {
				p.fQueue <- (<-p.fMsgPull.fQueue)
			}
			closed <- false
		}
	}()

	if <-closed {
		queue := make(chan message.IMessage)
		close(queue)
		return queue
	}

	p.fMutex.Lock()
	queue := p.fQueue
	p.fMutex.Unlock()

	return queue
}

func (p *sMessageQueue) newPseudoMessage() message.IMessage {
	msg, err := p.GetClient().EncryptPayload(
		p.GetClient().GetPubKey(),
		payload.NewPayload(0, []byte{1}),
	)
	if err != nil {
		panic(err)
	}
	return msg
}

func (p *sMessageQueue) readSignal() <-chan struct{} {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fMsgPull.fSignal
}
