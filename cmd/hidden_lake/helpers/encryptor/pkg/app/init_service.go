package app

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/internal/handler"
	hle_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/settings"
	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/message"
)

func (p *sApp) initServiceHTTP() {
	mux := http.NewServeMux()

	client := client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: p.fConfig.GetSettings().GetMessageSizeBytes(),
			FKeySizeBits:      p.fPrivKey.GetSize(),
		}),
		p.fPrivKey,
	)

	mux.HandleFunc(hle_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(hle_settings.CHandleEncryptPath, handler.HandleEncryptAPI(p.fConfig, p.fHTTPLogger, client))
	mux.HandleFunc(hle_settings.CHandleDecryptPath, handler.HandleDecryptAPI(p.fConfig, p.fHTTPLogger, client))

	p.fServiceHTTP = &http.Server{
		Addr:    p.fConfig.GetAddress().GetHTTP(),
		Handler: mux,
	}
}
