package queue

import (
	"context"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/utils"
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
		return utils.MergeErrors(ErrRunning, err)
	}
	defer func() { _ = p.fState.Disable(nil) }()

	for {
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		default:
			if err := p.fillMessagePool(pCtx); err != nil {
				return err
			}
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
	if p.queueHasLimit() {
		return ErrQueueLimit
	}
	p.fQueue <- net_message.NewMessage(
		p.fMsgSettings,
		payload.NewPayload(
			p.fNetworkMask,
			pMsg.ToBytes(),
		),
		p.fSettings.GetParallel(),
	)
	return nil
}

func (p *sMessageQueue) DequeueMessage(pCtx context.Context) net_message.IMessage {
	select {
	case <-pCtx.Done():
		return nil
	case <-time.After(p.fSettings.GetDuration()):
		select {
		case x := <-p.fQueue:
			// the main queue is checked first
			return x
		default:
			// take an existing message from any ready queue
			select {
			case <-pCtx.Done():
				return nil
			case x := <-p.fQueue:
				return x
			case x := <-p.fMsgPool.fQueue:
				return x
			}
		}
	}
}

func (p *sMessageQueue) fillMessagePool(pCtx context.Context) error {
	if p.poolHasLimit() {
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		case <-time.After(p.fSettings.GetDuration() / 2):
			return nil
		}
	}

	p.fMutex.Lock()
	msgSettings := p.fMsgSettings
	networkMask := p.fNetworkMask
	p.fMutex.Unlock()

	msg, err := p.fClient.EncryptPayload(
		p.fMsgPool.fReceiver,
		payload.NewPayload(0, []byte{1}),
	)
	if err != nil {
		return err
	}

	chNetMsg := make(chan net_message.IMessage)
	go func() {
		chNetMsg <- net_message.NewMessage(
			msgSettings,
			payload.NewPayload(
				networkMask,
				msg.ToBytes(),
			),
			p.fSettings.GetParallel(),
		)
	}()

	select {
	case <-pCtx.Done():
		return pCtx.Err()
	case x := <-chNetMsg:
		p.fMutex.Lock()
		settingsChanged := networkMask != p.fNetworkMask ||
			msgSettings.GetNetworkKey() != p.fMsgSettings.GetNetworkKey() ||
			msgSettings.GetWorkSizeBits() != p.fMsgSettings.GetWorkSizeBits()
		p.fMutex.Unlock()

		if !settingsChanged {
			p.fMsgPool.fQueue <- x
		}
		return nil
	}
}

func (p *sMessageQueue) queueHasLimit() bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	currLen := len(p.fQueue)
	return uint64(currLen) >= p.fSettings.GetMainCapacity()
}

func (p *sMessageQueue) poolHasLimit() bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	currLen := len(p.fMsgPool.fQueue)
	return uint64(currLen) >= p.fSettings.GetPoolCapacity()
}
