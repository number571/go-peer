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
	_ IQBTaskProcessor = &sQBTaskProcessor{}
)

type sQBTaskProcessor struct {
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

func NewQBTaskProcessor(
	pSettings ISettings,
	pVSettings IVSettings,
	pClient client.IClient,
	pReceiver asymmetric.IPubKey,
) IQBTaskProcessor {
	if pClient.GetPubKey().GetSize() != pReceiver.GetSize() {
		panic("pClient.GetPubKey().GetSize() != pReceiver.GetSize()")
	}
	return &sQBTaskProcessor{
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

func (p *sQBTaskProcessor) GetSettings() ISettings {
	return p.fSettings
}

func (p *sQBTaskProcessor) GetVSettings() IVSettings {
	return p.getVSettings()
}

func (p *sQBTaskProcessor) GetClient() client.IClient {
	return p.fClient
}

func (p *sQBTaskProcessor) Run(pCtx context.Context) error {
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

func (p *sQBTaskProcessor) runRandPoolFiller(pCtx context.Context, pWg *sync.WaitGroup, chErr chan<- error) {
	defer pWg.Done()

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

func (p *sQBTaskProcessor) runMainPoolFiller(pCtx context.Context, pWg *sync.WaitGroup, chErr chan<- error) {
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

func (p *sQBTaskProcessor) SetVSettings(pVSettings IVSettings) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fVSettings = pVSettings

	// clear all old queue state
	// not clear fMainPool.RawQueue
	for len(p.fMainPool.fQueue) > 0 {
		atomic.AddInt64(&p.fMainPool.fCount, -1)
		<-p.fMainPool.fQueue
	}

	for len(p.fRandPool.fQueue) > 0 {
		atomic.AddInt64(&p.fRandPool.fCount, -1)
		<-p.fRandPool.fQueue
	}
}

func (p *sQBTaskProcessor) EnqueueMessage(pPubKey asymmetric.IPubKey, pBytes []byte) error {
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

func (p *sQBTaskProcessor) DequeueMessage(pCtx context.Context) net_message.IMessage {
	queuePeriod := p.fSettings.GetQueuePeriod()
	addRandPeriod := time.Duration(random.NewCSPRNG().GetUint64() % uint64(p.fSettings.GetRandQueuePeriod()+1))

	for {
		select {
		case <-pCtx.Done():
			return nil
		case <-time.After(queuePeriod + addRandPeriod):
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
}

func (p *sQBTaskProcessor) fillMainPool(pCtx context.Context, pMsg []byte) error {
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

func (p *sQBTaskProcessor) fillRandPool(pCtx context.Context) error {
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

func (p *sQBTaskProcessor) getVSettings() IVSettings {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return p.fVSettings
}

func (p *sQBTaskProcessor) vSettingsNotChanged(oldVSettings IVSettings) bool {
	currVSettings := p.getVSettings()
	return currVSettings.GetNetworkKey() == oldVSettings.GetNetworkKey()
}
