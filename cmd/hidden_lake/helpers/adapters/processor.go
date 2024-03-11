package adapters

import (
	"context"
	"fmt"
	"io"
	"net/http"

	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/traffic/pkg/client"
	"github.com/number571/go-peer/internal/api"
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
		Addr:    pIncomingAddr,
		Handler: mux,
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
) error {
	for {
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		default:
		}
		msg, err := pConsumer.Consume(pCtx)
		if err != nil {
			pLogger.PushWarn(err.Error())
			continue
		}
		if msg == nil {
			pLogger.PushInfo("no new messages")
			continue
		}
		if err := pHltClient.PutMessage(msg); err != nil {
			pLogger.PushWarn(err.Error())
			continue
		}
		pLogger.PushInfo(fmt.Sprintf("message %X consumed", msg.GetHash()))
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
			api.Response(w, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		// get message from HLT
		msgStrAsBytes, err := io.ReadAll(r.Body)
		if err != nil {
			pLogger.PushWarn(err.Error())
			api.Response(w, http.StatusConflict, "failed: read body")
			return
		}

		msg, err := net_message.LoadMessage(pSettings, string(msgStrAsBytes))
		if err != nil {
			pLogger.PushWarn(err.Error())
			api.Response(w, http.StatusConflict, "failed: read message")
			return
		}

		if err := pProducer.Produce(pCtx, msg); err != nil {
			pLogger.PushWarn(err.Error())
			api.Response(w, http.StatusConflict, "failed: produce message")
			return
		}

		pLogger.PushInfo(fmt.Sprintf("message %X produced", msg.GetHash()))
		api.Response(w, http.StatusOK, "success: produce message")
	}
}
