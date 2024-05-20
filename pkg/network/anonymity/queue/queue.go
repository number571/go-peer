package queue

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/utils"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

var (
	_ IMessageQueue = &sMessageQueue{}
)

type sMessageQueue struct {
	fMutex sync.RWMutex
	fState state.IState

	fSettings  ISettings
	fVSettings IVSettings
	fClient    client.IClient

	fMainPool *sMainPool
	fVoidPool *sVoidPool
}

type sMainPool struct {
	fCount    int64 // atomic variable
	fQueue    chan net_message.IMessage
	fRawQueue chan []byte
}

type sVoidPool struct {
	fCount    int64 // atomic variable
	fQueue    chan net_message.IMessage
	fReceiver asymmetric.IPubKey
}

func NewMessageQueue(
	pSettings ISettings,
	pVSettings IVSettings,
	pClient client.IClient,
) IMessageQueue {
	return &sMessageQueue{
		fState:     state.NewBoolState(),
		fSettings:  pSettings,
		fVSettings: pVSettings,
		fClient:    pClient,
		fMainPool: &sMainPool{
			fQueue:    make(chan net_message.IMessage, pSettings.GetMainCapacity()),
			fRawQueue: make(chan []byte, pSettings.GetMainCapacity()),
		},
		fVoidPool: &sVoidPool{
			fQueue:    make(chan net_message.IMessage, pSettings.GetVoidCapacity()),
			fReceiver: asymmetric.NewRSAPrivKey(pClient.GetPrivKey().GetSize()).GetPubKey(),
		},
	}
}

func (p *sMessageQueue) GetSettings() ISettings {
	return p.fSettings
}

func (p *sMessageQueue) GetVSettings() IVSettings {
	return p.getVSettings()
}

func (p *sMessageQueue) GetClient() client.IClient {
	return p.fClient
}

func (p *sMessageQueue) Run(pCtx context.Context) error {
	if err := p.fState.Enable(nil); err != nil {
		return utils.MergeErrors(ErrRunning, err)
	}
	defer func() { _ = p.fState.Disable(nil) }()

	const numProcs = 2
	chBufErr := make(chan error, numProcs)

	wg := sync.WaitGroup{}
	wg.Add(numProcs)

	go p.runVoidPoolFiller(pCtx, &wg, chBufErr)
	go p.runMainPoolFiller(pCtx, &wg, chBufErr)

	wg.Wait()
	close(chBufErr)

	errList := make([]error, 0, numProcs)
	for err := range chBufErr {
		errList = append(errList, err)
	}
	return utils.MergeErrors(errList...)
}

func (p *sMessageQueue) runVoidPoolFiller(pCtx context.Context, pWg *sync.WaitGroup, chErr chan<- error) {
	defer pWg.Done()
	for {
		select {
		case <-pCtx.Done():
			chErr <- pCtx.Err()
			return
		default:
			if err := p.fillVoidPool(pCtx); err != nil {
				chErr <- err
				return
			}
		}
	}
}

func (p *sMessageQueue) runMainPoolFiller(pCtx context.Context, pWg *sync.WaitGroup, chErr chan<- error) {
	defer pWg.Done()
	for {
		select {
		case <-pCtx.Done():
			chErr <- pCtx.Err()
			return
		case x := <-p.fMainPool.fRawQueue:
			if err := p.fillMainPool(pCtx, x); err != nil {
				chErr <- err
				return
			}
		}
	}
}

func (p *sMessageQueue) SetVSettings(pVSettings IVSettings) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fVSettings = pVSettings

	// clear all old queue state
	// not clear fMainPool.RawQueue
	for len(p.fMainPool.fQueue) > 0 {
		atomic.AddInt64(&p.fMainPool.fCount, -1)
		<-p.fMainPool.fQueue
	}
	for len(p.fVoidPool.fQueue) > 0 {
		atomic.AddInt64(&p.fVoidPool.fCount, -1)
		<-p.fVoidPool.fQueue
	}
}

func (p *sMessageQueue) EnqueueMessage(pMsg []byte) error {
	incCount := atomic.AddInt64(&p.fMainPool.fCount, 1)
	if uint64(incCount) > p.fSettings.GetMainCapacity() {
		atomic.AddInt64(&p.fMainPool.fCount, -1)
		return ErrQueueLimit
	}
	p.fMainPool.fRawQueue <- pMsg
	return nil
}

func (p *sMessageQueue) DequeueMessage(pCtx context.Context) net_message.IMessage {
	randDuration := time.Duration(
		random.NewCSPRNG().GetUint64() % uint64(p.fSettings.GetRandDuration()+1),
	)

	select {
	case <-pCtx.Done():
		return nil
	case <-time.After(p.fSettings.GetDuration() + randDuration):
		select {
		case x := <-p.fMainPool.fQueue:
			// the main queue is checked first
			atomic.AddInt64(&p.fMainPool.fCount, -1)
			return x
		default:
			// take an existing message from any ready queue
			select {
			case <-pCtx.Done():
				return nil
			case x := <-p.fMainPool.fQueue:
				atomic.AddInt64(&p.fMainPool.fCount, -1)
				return x
			case x := <-p.fVoidPool.fQueue:
				atomic.AddInt64(&p.fVoidPool.fCount, -1)
				return x
			}
		}
	}
}

func (p *sMessageQueue) fillMainPool(pCtx context.Context, pMsg []byte) error {
	oldVSettings := p.getVSettings()
	chNetMsg := make(chan net_message.IMessage)

	go func() {
		chNetMsg <- net_message.NewMessage(
			net_message.NewSettings(&net_message.SSettings{
				FWorkSizeBits:       p.fSettings.GetWorkSizeBits(),
				FNetworkKey:         oldVSettings.GetNetworkKey(),
				FParallel:           p.fSettings.GetParallel(),
				FLimitVoidSizeBytes: p.fSettings.GetLimitVoidSizeBytes(),
			}),
			payload.NewPayload64(p.fSettings.GetNetworkMask(), pMsg),
		)
	}()

	select {
	case <-pCtx.Done():
		return pCtx.Err()

	case netMsg := <-chNetMsg:
		if p.vSettingsNotChanged(oldVSettings) {
			p.fMainPool.fQueue <- netMsg
		}
		return nil
	}
}

func (p *sMessageQueue) fillVoidPool(pCtx context.Context) error {
	incCount := atomic.AddInt64(&p.fVoidPool.fCount, 1)
	if uint64(incCount) > p.fSettings.GetVoidCapacity() {
		atomic.AddInt64(&p.fVoidPool.fCount, -1)
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		case <-time.After(p.fSettings.GetDuration() / 2):
			return nil
		}
	}

	msg, err := p.fClient.EncryptMessage(
		p.fVoidPool.fReceiver,
		payload.NewPayload64(0, []byte{1}).ToBytes(),
	)
	if err != nil {
		panic(err)
	}

	oldVSettings := p.getVSettings()
	chNetMsg := make(chan net_message.IMessage)
	go func() {
		chNetMsg <- net_message.NewMessage(
			net_message.NewSettings(&net_message.SSettings{
				FWorkSizeBits:       p.fSettings.GetWorkSizeBits(),
				FNetworkKey:         oldVSettings.GetNetworkKey(),
				FParallel:           p.fSettings.GetParallel(),
				FLimitVoidSizeBytes: p.fSettings.GetLimitVoidSizeBytes(),
			}),
			payload.NewPayload64(p.fSettings.GetNetworkMask(), msg),
		)
	}()

	select {
	case <-pCtx.Done():
		return pCtx.Err()

	case netMsg := <-chNetMsg:
		if p.vSettingsNotChanged(oldVSettings) {
			p.fVoidPool.fQueue <- netMsg
		}
		return nil
	}
}

func (p *sMessageQueue) getVSettings() IVSettings {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return p.fVSettings
}

func (p *sMessageQueue) vSettingsNotChanged(oldVSettings IVSettings) bool {
	currVSettings := p.getVSettings()
	return currVSettings.GetNetworkKey() == oldVSettings.GetNetworkKey()
}
