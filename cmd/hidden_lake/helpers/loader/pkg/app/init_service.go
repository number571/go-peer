package app

import (
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/internal/handler"
	hll_settings "github.com/number571/go-peer/cmd/hidden_lake/helpers/loader/pkg/settings"
)

func (p *sApp) initServiceHTTP() {
	mux := http.NewServeMux()

	mux.HandleFunc(hll_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(hll_settings.CHandleNetworkTransferPath, handler.HandleNetworkTransferAPI(p.fConfig, p.fHTTPLogger))
	mux.HandleFunc(hll_settings.CHandleConfigSettings, handler.HandleConfigSettingsAPI(p.fConfig, p.fHTTPLogger))

	p.fServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetHTTP(),
		ReadTimeout: (5 * time.Second),
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
	}
}