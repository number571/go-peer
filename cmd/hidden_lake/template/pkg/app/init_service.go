package app

import (
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/template/internal/handler"
	hl_t_settings "github.com/number571/go-peer/cmd/hidden_lake/template/pkg/settings"
)

func (p *sApp) initServiceHTTP() {
	mux := http.NewServeMux()

	mux.HandleFunc(hl_t_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(hl_t_settings.CHandleConfigSettingsPath, handler.HandleConfigSettingsAPI(p.fConfig, p.fHTTPLogger))

	p.fServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetHTTP(),
		ReadTimeout: (5 * time.Second),
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
	}
}
