package app

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/config"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network/anonymity"

	hlt_client "github.com/number571/go-peer/cmd/hidden_lake/traffic/pkg/client"
)

func (p *sApp) loadTrafficMessages() {
	cfg := p.fWrapper.GetConfig()
	for _, conn := range cfg.GetBackupConnections() {
		go handleTrafficMessages(cfg, p.fNode, conn)
	}
}

func handleTrafficMessages(pCfg config.IConfig, pNode anonymity.INode, pConn string) {
	sett := pCfg.GetSettings()
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
		if err := pNode.HandleMessage(msg); err != nil {
			continue
		}
	}
}
