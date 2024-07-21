package app

import (
	"context"
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/remoter/internal/handler"
	hlr_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/remoter/pkg/settings"
)

func (p *sApp) initIncomingServiceHTTP(pCtx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		hlr_settings.CExecPath,
		handler.HandleIncomigExecHTTP(pCtx, p.fConfig, p.fHTTPLogger),
	) // POST

	p.fIncServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetIncoming(),
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: (5 * time.Second),
	}
}
