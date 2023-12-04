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
	fIsRun bool
	fMutex sync.Mutex

	fNetworkMask uint64
	fMsgSettings net_message.ISettings

	fSettings ISettings
	fClient   client.IClient
	fQueue    chan net_message.IMessage
	fMsgPool  sPool
}

type sPool struct {
	fSignal   chan struct{}
	fQueue    chan net_message.IMessage
	fReceiver asymmetric.IPubKey
}

func NewMessageQueue(pSett ISettings, pClient client.IClient) IMessageQueue {
	return &sMessageQueue{
		fMsgSettings: net_message.NewSettings(&net_message.SSettings{}),
		fSettings:    pSett,
		fClient:      pClient,
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
				if p.poolHasLimit() {
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

func (p *sMessageQueue) WithNetworkSettings(pNetworkMask uint64, pMsgSettings net_message.ISettings) IMessageQueue {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	for len(p.fQueue) > 0 {
		<-p.fQueue
	}
	for len(p.fMsgPool.fQueue) > 0 {
		<-p.fMsgPool.fQueue
	}

	p.fNetworkMask = pNetworkMask
	p.fMsgSettings = pMsgSettings

	return p
}

func (p *sMessageQueue) EnqueueMessage(pMsg message.IMessage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if uint64(len(p.fQueue)) >= p.fSettings.GetMainCapacity() {
		return errors.New("queue already full, need wait and retry")
	}

	p.fQueue <- net_message.NewMessage(
		p.fMsgSettings,
		payload.NewPayload(
			p.fNetworkMask,
			pMsg.ToBytes(),
		),
	)
	return nil
}

func (p *sMessageQueue) DequeueMessage() net_message.IMessage {
	var (
		result net_message.IMessage
		closed = make(chan bool)
	)

	go func() {
		select {
		case <-p.readSignal():
			closed <- true
			return
		case <-time.After(p.fSettings.GetDuration()):
			defer func() { closed <- false }()

			p.fMutex.Lock()
			queueLen := len(p.fQueue)
			p.fMutex.Unlock()

			if queueLen == 0 {
				result = <-p.fMsgPool.fQueue
				return
			}

			result = <-p.fQueue
			return
		}
	}()

	<-closed
	return result
}

func (p *sMessageQueue) newPseudoNetworkMessage() net_message.IMessage {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	msg, err := p.fClient.EncryptPayload(
		p.fMsgPool.fReceiver,
		payload.NewPayload(0, []byte{1}),
	)
	if err != nil {
		panic(err)
	}

	return net_message.NewMessage(
		p.fMsgSettings,
		payload.NewPayload(
			p.fNetworkMask,
			msg.ToBytes(),
		),
	)
}

func (p *sMessageQueue) readSignal() <-chan struct{} {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.fMsgPool.fSignal
}

func (p *sMessageQueue) poolHasLimit() bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	currLen := len(p.fMsgPool.fQueue)
	return uint64(currLen) >= p.fSettings.GetPoolCapacity()
}
