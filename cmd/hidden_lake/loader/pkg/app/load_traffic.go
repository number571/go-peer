package app

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/loader/internal/config"
	"github.com/number571/go-peer/pkg/encoding"

	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

func (p *sApp) transferMessages(pCtx context.Context) {
	cfg := p.fConfig

	producerClients := make([]hlt_client.IClient, 0, len(cfg.GetConsumers()))
	for _, p := range cfg.GetConsumers() {
		producerClients = append(producerClients, makeHLTClient(cfg, p))
	}

	for _, p := range producerClients {
		go transferToConsumers(pCtx, cfg, p)
	}
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
			for _, c := range consumerClients {
				go func(c hlt_client.IClient) {
					defer wg.Done()
					_ = c.PutMessage(msg)
				}(c)
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
			&http.Client{Timeout: time.Minute},
			net_message.NewSettings(&net_message.SSettings{
				FWorkSizeBits: sett.GetWorkSizeBits(),
				FNetworkKey:   sett.GetNetworkKey(),
			}),
		),
	)
}
