package app

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"

	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
)

func (p *sApp) loadTrafficMessages() {
	cfg := p.fConfig
	for _, conn := range cfg.GetBackupConnections() {
		go handleTrafficMessages(cfg, conn)
	}
}

func handleTrafficMessages(pCfg config.IConfig, pConn string) {
	hlsClient := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			fmt.Sprintf("http://%s", pCfg.GetConnection()),
			&http.Client{Timeout: time.Minute},
		),
	)

	sett, err := hlsClient.GetSettings()
	if err != nil {
		return
	}

	hltClient := hlt_client.NewClient(
		hlt_client.NewBuilder(),
		hlt_client.NewRequester(
			fmt.Sprintf("http://%s", pConn),
			&http.Client{Timeout: time.Minute},
			message.NewSettings(&message.SSettings{
				FWorkSizeBits:     sett.GetWorkSizeBits(),
				FMessageSizeBytes: sett.GetMessageSizeBytes(),
			}),
		),
	)

	hashes, err := hltClient.GetHashes()
	if err != nil {
		return
	}

	for i, hash := range hashes {
		if uint64(i) >= pCfg.GetSettings().GetMessagesCapacity() {
			break
		}
		msg, err := hltClient.GetMessage(hash)
		if err != nil {
			continue
		}
		bytesHash := encoding.HexDecode(hash)
		if !bytes.Equal(msg.GetBody().GetHash(), bytesHash) {
			break
		}
		if err := hlsClient.HandleMessage(msg); err != nil {
			continue
		}
	}
}
