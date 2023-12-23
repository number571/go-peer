package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/state"
)

var (
	_ IMessageQueue = &sMessageQueue{}
)

type sMessageQueue struct {
	fState state.IState
	fMutex sync.Mutex

	fNetworkMask uint64
	fMsgSettings net_message.ISettings

	fSettings ISettings
	fClient   client.IClient
	fQueue    chan net_message.IMessage
	fMsgPool  sPool
}

type sPool struct {
	fQueue    chan net_message.IMessage
	fReceiver asymmetric.IPubKey
}

func NewMessageQueue(pSett ISettings, pClient client.IClient) IMessageQueue {
	return &sMessageQueue{
		fState:       state.NewBoolState(),
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

func (p *sMessageQueue) Run(pCtx context.Context) error {
	if err := p.fState.Enable(nil); err != nil {
		return fmt.Errorf("queue running error: %w", err)
	}
	defer func() {
		if err := p.fState.Disable(nil); err != nil {
			panic(err)
		}
	}()

	for {
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		case <-time.After(p.fSettings.GetDuration() / 2):
			if p.poolHasLimit() {
				continue
			}
			p.fMsgPool.fQueue <- p.newPseudoNetworkMessage()
		}
	}
}

func (p *sMessageQueue) WithNetworkSettings(pNetworkMask uint64, pMsgSettings net_message.ISettings) IMessageQueue {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fNetworkMask = pNetworkMask
	p.fMsgSettings = pMsgSettings

	for len(p.fQueue) > 0 {
		<-p.fQueue
	}
	for len(p.fMsgPool.fQueue) > 0 {
		<-p.fMsgPool.fQueue
	}

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

func (p *sMessageQueue) DequeueMessage(pCtx context.Context) net_message.IMessage {
	select {
	case <-pCtx.Done():
		return nil
	case <-time.After(p.fSettings.GetDuration()):
		p.fMutex.Lock()
		queueLen := len(p.fQueue)
		p.fMutex.Unlock()

		if queueLen == 0 {
			return <-p.fMsgPool.fQueue
		}
		return <-p.fQueue
	}
}

func (p *sMessageQueue) newPseudoNetworkMessage() net_message.IMessage {
	msg, err := p.fClient.EncryptPayload(
		p.fMsgPool.fReceiver,
		payload.NewPayload(0, []byte{1}),
	)
	if err != nil {
		panic(err)
	}

	p.fMutex.Lock()
	msgSettings := p.fMsgSettings
	networkMask := p.fNetworkMask
	p.fMutex.Unlock()

	return net_message.NewMessage(
		msgSettings,
		payload.NewPayload(
			networkMask,
			msg.ToBytes(),
		),
	)
}

func (p *sMessageQueue) poolHasLimit() bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	currLen := len(p.fMsgPool.fQueue)
	return uint64(currLen) >= p.fSettings.GetPoolCapacity()
}
