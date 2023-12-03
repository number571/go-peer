package queue

import (
	"errors"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
)

var (
	_ IMessageQueue = &sMessageQueue{}
)

type sMessageQueue struct {
	fIsRun       bool
	fMutex       sync.Mutex
	fSettings    ISettings
	fClient      client.IClient
	fNetSettFunc INetworkSettingsFunc
	fQueue       chan net_message.IMessage
	fMsgPool     sPool
}

type sPool struct {
	fSignal   chan struct{}
	fQueue    chan net_message.IMessage
	fReceiver asymmetric.IPubKey
}

func NewMessageQueue(pSett ISettings, pClient client.IClient, pNetSettFunc INetworkSettingsFunc) IMessageQueue {
	return &sMessageQueue{
		fSettings:    pSett,
		fClient:      pClient,
		fNetSettFunc: pNetSettFunc,
		fQueue:       make(chan net_message.IMessage, pSett.GetMainCapacity()),
		fMsgPool: sPool{
			fQueue:    make(chan net_message.IMessage, pSett.GetPoolCapacity()),
			fReceiver: asymmetric.NewRSAPrivKey(pClient.GetPrivKey().GetSize()).GetPubKey(),
		},
	}
}

func (p *sMessageQueue) GetSettings() ISettings {
	return p.fSettings
}

func (p *sMessageQueue) GetClient() client.IClient {
	return p.fClient
}

func (p *sMessageQueue) ClearQueue() {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fQueue = make(chan net_message.IMessage, p.fSettings.GetMainCapacity())
	p.fMsgPool.fQueue = make(chan net_message.IMessage, p.fSettings.GetPoolCapacity())
	p.fMsgPool.fReceiver = asymmetric.NewRSAPrivKey(p.fClient.GetPrivKey().GetSize()).GetPubKey()
}

func (p *sMessageQueue) Run() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fIsRun {
		return errors.New("queue already running")
	}
	p.fIsRun = true

	p.fMsgPool.fSignal = make(chan struct{})
	go func() {
		for {
			select {
			case <-p.readSignal():
				return
			case <-time.After(p.fSettings.GetDuration() / 2):
				if p.hasLimit() {
					continue
				}
				p.fMsgPool.fQueue <- p.newPseudoNetworkMessage()
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

	close(p.fMsgPool.fSignal)
	return nil
}

func (p *sMessageQueue) EnqueueMessage(pMsg message.IMessage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if uint64(len(p.fQueue)) >= p.fSettings.GetMainCapacity() {
		return errors.New("queue already full, need wait and retry")
	}

	networkMask, messageSettings := p.fNetSettFunc()
	p.fQueue <- net_message.NewMessage(
		messageSettings,
		payload.NewPayload(
			networkMask,
			pMsg.ToBytes(),
		),
	)
	return nil
}

func (p *sMessageQueue) DequeueMessage() <-chan net_message.IMessage {
	closed := make(chan bool)

	go func() {
		select {
		case <-p.readSignal():
			closed <- true
			return
		case <-time.After(p.fSettings.GetDuration()):
			p.fMutex.Lock()
			defer p.fMutex.Unlock()

			if len(p.fQueue) == 0 {
				p.fQueue <- (<-p.fMsgPool.fQueue)
			}
			closed <- false
		}
	}()

	if <-closed {
		queue := make(chan net_message.IMessage)
		close(queue)
		return queue
	}

	p.fMutex.Lock()
	queue := p.fQueue
	p.fMutex.Unlock()

	return queue
}

func (p *sMessageQueue) newPseudoNetworkMessage() net_message.IMessage {
	msg, err := p.fClient.EncryptPayload(
		p.fMsgPool.fReceiver,
		payload.NewPayload(0, []byte{1}),
	)
	if err != nil {
		panic(err)
	}
	networkMask, messageSettings := p.fNetSettFunc()
	return net_message.NewMessage(
		messageSettings,
		payload.NewPayload(
			networkMask,
			msg.ToBytes(),
		),
	)
}

func (p *sMessageQueue) readSignal() <-chan struct{} {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fMsgPool.fSignal
}

func (p *sMessageQueue) hasLimit() bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	currLen := len(p.fMsgPool.fQueue)
	return uint64(currLen) >= p.fSettings.GetPoolCapacity()
}
