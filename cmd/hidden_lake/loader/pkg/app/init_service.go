package app

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/loader/internal/handler"
	hll_settings "github.com/number571/go-peer/cmd/hidden_lake/loader/pkg/settings"
)

func (p *sApp) initServiceHTTP() {
	mux := http.NewServeMux()

	mux.HandleFunc(hll_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(hll_settings.CHandleTransferPath, handler.HandleTransferAPI(p.fConfig, p.fHTTPLogger))

	p.fServiceHTTP = &http.Server{
		Addr:    p.fConfig.GetAddress().GetHTTP(),
		Handler: mux,
	}
}
