package queue

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/state"
	"github.com/number571/go-peer/pkg/utils"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

var (
	_ IMessageQueueProcessor = &sMessageQueueProcessor{}
)

type sMessageQueueProcessor struct {
	fMutex sync.RWMutex
	fState state.IState

	fSettings  ISettings
	fVSettings IVSettings
	fClient    client.IClient

	fMainPool *sMainPool
	fRandPool *sRandPool
}

type sMainPool struct {
	fCount    int64 // atomic variable
	fQueue    chan net_message.IMessage
	fRawQueue chan []byte
}

type sRandPool struct {
	fCount    int64 // atomic variable
	fQueue    chan net_message.IMessage
	fReceiver asymmetric.IPubKey
}

func NewMessageQueueProcessor(
	pSettings ISettings,
	pVSettings IVSettings,
	pClient client.IClient,
	pReceiver asymmetric.IPubKey,
) IMessageQueueProcessor {
	if pSettings.GetQueuePeriod() != 0 {
		if pClient.GetPubKey().GetSize() != pReceiver.GetSize() {
			panic("pClient.GetPubKey().GetSize() != pReceiver.GetSize()")
		}
	}
	return &sMessageQueueProcessor{
		fState:     state.NewBoolState(),
		fSettings:  pSettings,
		fVSettings: pVSettings,
		fClient:    pClient,
		fMainPool: &sMainPool{
			fQueue:    make(chan net_message.IMessage, pSettings.GetMainPoolCapacity()),
			fRawQueue: make(chan []byte, pSettings.GetMainPoolCapacity()),
		},
		fRandPool: &sRandPool{
			fQueue:    make(chan net_message.IMessage, pSettings.GetRandPoolCapacity()),
			fReceiver: pReceiver,
		},
	}
}

func (p *sMessageQueueProcessor) GetSettings() ISettings {
	return p.fSettings
}

func (p *sMessageQueueProcessor) GetVSettings() IVSettings {
	return p.getVSettings()
}

func (p *sMessageQueueProcessor) GetClient() client.IClient {
	return p.fClient
}

func (p *sMessageQueueProcessor) Run(pCtx context.Context) error {
	if err := p.fState.Enable(nil); err != nil {
		return utils.MergeErrors(ErrRunning, err)
	}
	defer func() { _ = p.fState.Disable(nil) }()

	const numProcs = 2
	chBufErr := make(chan error, numProcs)

	wg := sync.WaitGroup{}
	wg.Add(numProcs)

	go p.runRandPoolFiller(pCtx, &wg, chBufErr)
	go p.runMainPoolFiller(pCtx, &wg, chBufErr)

	wg.Wait()
	close(chBufErr)

	errList := make([]error, 0, numProcs)
	for err := range chBufErr {
		errList = append(errList, err)
	}
	return utils.MergeErrors(errList...)
}

func (p *sMessageQueueProcessor) runRandPoolFiller(pCtx context.Context, pWg *sync.WaitGroup, chErr chan<- error) {
	defer pWg.Done()

	if p.fSettings.GetQueuePeriod() == 0 { // if QB=false
		<-pCtx.Done()
		chErr <- pCtx.Err()
		return
	}

	for {
		select {
		case <-pCtx.Done():
			chErr <- pCtx.Err()
			return
		default:
			if err := p.fillRandPool(pCtx); err != nil {
				chErr <- err
				return
			}
		}
	}
}

func (p *sMessageQueueProcessor) runMainPoolFiller(pCtx context.Context, pWg *sync.WaitGroup, chErr chan<- error) {
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

func (p *sMessageQueueProcessor) SetVSettings(pVSettings IVSettings) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fVSettings = pVSettings

	// clear all old queue state
	// not clear fMainPool.RawQueue
	for len(p.fMainPool.fQueue) > 0 {
		atomic.AddInt64(&p.fMainPool.fCount, -1)
		<-p.fMainPool.fQueue
	}

	if p.fSettings.GetQueuePeriod() != 0 { // if QB=true
		for len(p.fRandPool.fQueue) > 0 {
			atomic.AddInt64(&p.fRandPool.fCount, -1)
			<-p.fRandPool.fQueue
		}
	}
}

func (p *sMessageQueueProcessor) EnqueueMessage(pPubKey asymmetric.IPubKey, pBytes []byte) error {
	incCount := atomic.AddInt64(&p.fMainPool.fCount, 1)
	if uint64(incCount) > p.fSettings.GetMainPoolCapacity() {
		atomic.AddInt64(&p.fMainPool.fCount, -1)
		return ErrQueueLimit
	}
	rawMsg, err := p.fClient.EncryptMessage(pPubKey, pBytes)
	if err != nil {
		atomic.AddInt64(&p.fMainPool.fCount, -1)
		return utils.MergeErrors(ErrEncryptMessage, err)
	}
	p.fMainPool.fRawQueue <- rawMsg
	return nil
}

func (p *sMessageQueueProcessor) DequeueMessage(pCtx context.Context) net_message.IMessage {
	if p.fSettings.GetQueuePeriod() == 0 { // if QB=false
		select {
		case <-pCtx.Done():
			return nil
		case x := <-p.fMainPool.fQueue:
			atomic.AddInt64(&p.fMainPool.fCount, -1)
			return x
		}
	}

	randUint64 := random.NewCSPRNG().GetUint64()
	randQueuePeriod := time.Duration(randUint64 % uint64(p.fSettings.GetRandQueuePeriod()+1))

	select {
	case <-pCtx.Done():
		return nil
	case <-time.After(p.fSettings.GetQueuePeriod() + randQueuePeriod):
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
			case x := <-p.fRandPool.fQueue:
				atomic.AddInt64(&p.fRandPool.fCount, -1)
				return x
			}
		}
	}
}

func (p *sMessageQueueProcessor) fillMainPool(pCtx context.Context, pMsg []byte) error {
	oldVSettings := p.getVSettings()
	chNetMsg := make(chan net_message.IMessage)

	go func() {
		chNetMsg <- net_message.NewMessage(
			net_message.NewSettings(&net_message.SSettings{
				FWorkSizeBits:         p.fSettings.GetWorkSizeBits(),
				FNetworkKey:           oldVSettings.GetNetworkKey(),
				FParallel:             p.fSettings.GetParallel(),
				FRandMessageSizeBytes: p.fSettings.GetRandMessageSizeBytes(),
			}),
			payload.NewPayload32(p.fSettings.GetNetworkMask(), pMsg),
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

func (p *sMessageQueueProcessor) fillRandPool(pCtx context.Context) error {
	incCount := atomic.AddInt64(&p.fRandPool.fCount, 1)
	if uint64(incCount) > p.fSettings.GetRandPoolCapacity() {
		atomic.AddInt64(&p.fRandPool.fCount, -1)
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		case <-time.After(p.fSettings.GetQueuePeriod() / 2):
			return nil
		}
	}

	msg, err := p.fClient.EncryptMessage(
		p.fRandPool.fReceiver,
		random.NewCSPRNG().GetBytes(encoding.CSizeUint64),
	)
	if err != nil {
		panic(err)
	}

	oldVSettings := p.getVSettings()
	chNetMsg := make(chan net_message.IMessage)

	go func() {
		chNetMsg <- net_message.NewMessage(
			net_message.NewSettings(&net_message.SSettings{
				FWorkSizeBits:         p.fSettings.GetWorkSizeBits(),
				FNetworkKey:           oldVSettings.GetNetworkKey(),
				FParallel:             p.fSettings.GetParallel(),
				FRandMessageSizeBytes: p.fSettings.GetRandMessageSizeBytes(),
			}),
			payload.NewPayload32(p.fSettings.GetNetworkMask(), msg),
		)
	}()

	select {
	case <-pCtx.Done():
		return pCtx.Err()
	case netMsg := <-chNetMsg:
		if p.vSettingsNotChanged(oldVSettings) {
			p.fRandPool.fQueue <- netMsg
		}
		return nil
	}
}

func (p *sMessageQueueProcessor) getVSettings() IVSettings {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return p.fVSettings
}

func (p *sMessageQueueProcessor) vSettingsNotChanged(oldVSettings IVSettings) bool {
	currVSettings := p.getVSettings()
	return currVSettings.GetNetworkKey() == oldVSettings.GetNetworkKey()
}
