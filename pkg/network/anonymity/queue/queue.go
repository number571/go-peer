package queue

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/state"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

var (
	_ IQBProblemProcessor = &sQBProblemProcessor{}
)

type sQBProblemProcessor struct {
	fState state.IState

	fSettings ISettings
	fClient   client.IClient

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

func NewQBProblemProcessor(pSettings ISettings, pClient client.IClient) IQBProblemProcessor {
	poolCap := pSettings.GetPoolCapacity()
	return &sQBProblemProcessor{
		fState:    state.NewBoolState(),
		fSettings: pSettings,
		fClient:   pClient,
		fMainPool: &sMainPool{
			fQueue:    make(chan net_message.IMessage, poolCap[0]),
			fRawQueue: make(chan []byte, poolCap[0]),
		},
		fRandPool: &sRandPool{
			fQueue:    make(chan net_message.IMessage, poolCap[1]),
			fReceiver: asymmetric.NewPrivKey().GetPubKey(),
		},
	}
}

func (p *sQBProblemProcessor) GetSettings() ISettings {
	return p.fSettings
}

func (p *sQBProblemProcessor) GetClient() client.IClient {
	return p.fClient
}

func (p *sQBProblemProcessor) Run(pCtx context.Context) error {
	ctx, cancel := context.WithCancel(pCtx)
	defer cancel()

	if err := p.fState.Enable(nil); err != nil {
		return errors.Join(ErrRunning, err)
	}
	defer func() { _ = p.fState.Disable(nil) }()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go p.runRandPoolFiller(ctx, cancel, &wg)
	go p.runMainPoolFiller(ctx, cancel, &wg)

	wg.Wait()
	return ctx.Err()
}

func (p *sQBProblemProcessor) runRandPoolFiller(pCtx context.Context, pCancel func(), pWG *sync.WaitGroup) {
	defer func() {
		pWG.Done()
		pCancel()
	}()
	for {
		select {
		case <-pCtx.Done():
			return
		default:
			if err := p.fillRandPool(pCtx); err != nil {
				return
			}
		}
	}
}

func (p *sQBProblemProcessor) runMainPoolFiller(pCtx context.Context, pCancel func(), pWG *sync.WaitGroup) {
	defer func() {
		pWG.Done()
		pCancel()
	}()
	for {
		select {
		case <-pCtx.Done():
			return
		case msg := <-p.fMainPool.fRawQueue:
			if err := p.pushMessage(pCtx, p.fMainPool.fQueue, msg); err != nil {
				return
			}
		}
	}
}

func (p *sQBProblemProcessor) EnqueueMessage(pPubKey asymmetric.IPubKey, pBytes []byte) error {
	incCount := atomic.AddInt64(&p.fMainPool.fCount, 1)
	if uint64(incCount) > uint64(cap(p.fMainPool.fQueue)) {
		atomic.AddInt64(&p.fMainPool.fCount, -1)
		return ErrQueueLimit
	}
	rawMsg, err := p.fClient.EncryptMessage(pPubKey, pBytes)
	if err != nil {
		atomic.AddInt64(&p.fMainPool.fCount, -1)
		return errors.Join(ErrEncryptMessage, err)
	}
	p.fMainPool.fRawQueue <- rawMsg
	return nil
}

func (p *sQBProblemProcessor) DequeueMessage(pCtx context.Context) net_message.IMessage {
	for {
		select {
		case <-pCtx.Done():
			return nil
		case <-time.After(p.fSettings.GetQueuePeriod()):
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

func (p *sQBProblemProcessor) fillRandPool(pCtx context.Context) error {
	incCount := atomic.AddInt64(&p.fRandPool.fCount, 1)
	if uint64(incCount) > uint64(cap(p.fRandPool.fQueue)) {
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
		random.NewRandom().GetBytes(encoding.CSizeUint64),
	)
	if err != nil {
		panic(err)
	}
	return p.pushMessage(pCtx, p.fRandPool.fQueue, msg)
}

func (p *sQBProblemProcessor) pushMessage(pCtx context.Context, pQueue chan<- net_message.IMessage, pMsg []byte) error {
	chNetMsg := make(chan net_message.IMessage)
	go func() {
		chNetMsg <- net_message.NewMessage(
			p.fSettings.GetMessageConstructSettings(),
			payload.NewPayload32(p.fSettings.GetNetworkMask(), pMsg),
		)
	}()
	select {
	case <-pCtx.Done():
		return pCtx.Err()
	case netMsg := <-chNetMsg:
		pQueue <- netMsg
		return nil
	}
}
