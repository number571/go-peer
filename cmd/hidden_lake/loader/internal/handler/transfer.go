package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"

	"github.com/number571/go-peer/cmd/hidden_lake/loader/internal/config"
	hll_settings "github.com/number571/go-peer/cmd/hidden_lake/loader/pkg/settings"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

type sTransfer struct {
	fMutex  sync.Mutex
	fIsRun  bool
	fConfig config.IConfig
	fCancel context.CancelFunc
}

func HandleTransferAPI(pConfig config.IConfig, pLogger logger.ILogger) http.HandlerFunc {
	transfer := &sTransfer{fConfig: pConfig}

	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hll_settings.CServiceName, pR)

		if pR.Method != http.MethodPost && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		switch pR.Method {
		case http.MethodPost:
			if err := transfer.run(); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("already_running"))
				api.Response(pW, http.StatusOK, "failed: already running")
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, "success: run transfer")
			return

		case http.MethodDelete:
			if err := transfer.stop(); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("already_stopped"))
				api.Response(pW, http.StatusOK, "failed: already atopped")
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			api.Response(pW, http.StatusOK, "success: stop transfer")
			return
		}
	}
}

func (p *sTransfer) run() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if p.fIsRun {
		return errors.New("loader already is running")
	}
	p.fIsRun = true

	ctx := context.Background()
	ctxWithCancel, cancelFunction := context.WithCancel(ctx)

	p.fCancel = cancelFunction
	go p.transferMessages(ctxWithCancel)

	return nil
}

func (p *sTransfer) stop() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if !p.fIsRun {
		return errors.New("loader already is stopped")
	}
	p.fIsRun = false

	p.fCancel()
	return nil
}

func (p *sTransfer) transferMessages(pCtx context.Context) {
	cfg := p.fConfig

	producerClients := make([]hlt_client.IClient, 0, len(cfg.GetProducers()))
	for _, producer := range cfg.GetProducers() {
		producerClients = append(producerClients, makeHLTClient(cfg, producer))
	}

	wg := sync.WaitGroup{}
	wg.Add(len(producerClients))

	for _, producer := range producerClients {
		go func(p hlt_client.IClient) {
			defer wg.Done()
			transferToConsumers(pCtx, cfg, p)
		}(producer)
	}

	wg.Wait()
	_ = p.stop()
}

func transferToConsumers(pCtx context.Context, pCfg config.IConfig, pProducer hlt_client.IClient) {
	consumerClients := make([]hlt_client.IClient, 0, len(pCfg.GetConsumers()))
	for _, c := range pCfg.GetConsumers() {
		consumerClients = append(consumerClients, makeHLTClient(pCfg, c))
	}

	hashes, err := pProducer.GetHashes()
	if err != nil {
		return
	}

	for i, hash := range hashes {
		select {
		case <-pCtx.Done():
			return
		default:
			if uint64(i) >= pCfg.GetSettings().GetMessagesCapacity() {
				break
			}

			msg, err := pProducer.GetMessage(hash)
			if err != nil {
				continue
			}

			bytesHash := encoding.HexDecode(hash)
			if !bytes.Equal(msg.GetHash(), bytesHash) {
				break
			}

			wg := sync.WaitGroup{}
			wg.Add(len(consumerClients))
			for _, consumer := range consumerClients {
				go func(c hlt_client.IClient) {
					defer wg.Done()
					_ = c.PutMessage(msg)
				}(consumer)
			}
			wg.Wait()
		}
	}
}

func makeHLTClient(pCfg config.IConfig, pConn string) hlt_client.IClient {
	sett := pCfg.GetSettings()
	return hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			fmt.Sprintf("http://%s", pConn),
			&http.Client{Timeout: time.Minute / 2},
			net_message.NewSettings(&net_message.SSettings{
				FWorkSizeBits: sett.GetWorkSizeBits(),
				FNetworkKey:   sett.GetNetworkKey(),
			}),
		),
	)
}
