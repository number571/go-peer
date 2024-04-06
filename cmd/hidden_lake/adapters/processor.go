package adapters

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/client"
	"github.com/number571/go-peer/internal/api"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

func ProduceProcessor(
	pCtx context.Context,
	pProducer IAdaptedProducer,
	pLogger logger.ILogger,
	pSettings net_message.ISettings,
	pIncomingAddr string,
) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/adapter", produceHandler(pCtx, pProducer, pLogger, pSettings))

	server := &http.Server{
		Addr:              pIncomingAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-pCtx.Done()
		server.Close()
	}()

	return server.ListenAndServe()
}

func ConsumeProcessor(
	pCtx context.Context,
	pConsumer IAdaptedConsumer,
	pLogger logger.ILogger,
	pHltClient hlt_client.IClient,
	pWaitTimeout time.Duration,
) error {
	for {
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		default:
		}
		msg, err := pConsumer.Consume(pCtx)
		if err != nil {
			wait(pCtx, pWaitTimeout)
			pLogger.PushWarn(err.Error())
			continue
		}
		if msg == nil {
			wait(pCtx, pWaitTimeout)
			pLogger.PushInfo("no new messages")
			continue
		}
		if err := pHltClient.PutMessage(pCtx, msg); err != nil {
			wait(pCtx, pWaitTimeout)
			pLogger.PushWarn(err.Error())
			continue
		}
		pLogger.PushInfo(fmt.Sprintf("message %X consumed", msg.GetHash()))
	}
}

func wait(pCtx context.Context, pWaitTimeout time.Duration) {
	randDuration := time.Duration(
		random.NewStdPRNG().GetUint64() % uint64(pWaitTimeout+1),
	)
	select {
	case <-pCtx.Done():
	case <-time.After(pWaitTimeout + randDuration):
	}
}

func produceHandler(
	pCtx context.Context,
	pProducer IAdaptedProducer,
	pLogger logger.ILogger,
	pSettings net_message.ISettings,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			pLogger.PushWarn("got method != post")
			_ = api.Response(w, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		// get message from HLT
		msgStrAsBytes, err := io.ReadAll(r.Body)
		if err != nil {
			pLogger.PushWarn(err.Error())
			_ = api.Response(w, http.StatusConflict, "failed: read body")
			return
		}

		msg, err := net_message.LoadMessage(pSettings, string(msgStrAsBytes))
		if err != nil {
			pLogger.PushWarn(err.Error())
			_ = api.Response(w, http.StatusConflict, "failed: read message")
			return
		}

		if err := pProducer.Produce(pCtx, msg); err != nil {
			pLogger.PushWarn(err.Error())
			_ = api.Response(w, http.StatusConflict, "failed: produce message")
			return
		}

		pLogger.PushInfo(fmt.Sprintf("message %X produced", msg.GetHash()))
		_ = api.Response(w, http.StatusOK, "success: produce message")
	}
}
