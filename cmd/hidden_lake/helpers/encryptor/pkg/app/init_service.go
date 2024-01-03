package app

import (
	"net/http"
	"time"

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
	mux.HandleFunc(hle_settings.CHandleMessageEncryptPath, handler.HandleMessageEncryptAPI(p.fConfig, p.fHTTPLogger, client, p.fParallel))
	mux.HandleFunc(hle_settings.CHandleMessageDecryptPath, handler.HandleMessageDecryptAPI(p.fConfig, p.fHTTPLogger, client))
	mux.HandleFunc(hle_settings.CHandleServicePubKeyPath, handler.HandleServicePubKeyAPI(p.fHTTPLogger, client.GetPubKey()))
	mux.HandleFunc(hle_settings.CHandleConfigSettings, handler.HandleConfigSettingsAPI(p.fConfig, p.fHTTPLogger))

	p.fServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetHTTP(),
		ReadTimeout: (5 * time.Second),
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
	}
}
