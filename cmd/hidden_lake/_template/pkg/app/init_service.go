package app

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/_template/internal/handler"
	hl_t_settings "github.com/number571/go-peer/cmd/hidden_lake/_template/pkg/settings"
)

func (p *sApp) initServiceHTTP() {
	mux := http.NewServeMux()

	// TODO: need implementation
	mux.HandleFunc(hl_t_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))

	p.fServiceHTTP = &http.Server{
		Addr:    p.fConfig.GetAddress().GetHTTP(),
		Handler: mux,
	}
}
