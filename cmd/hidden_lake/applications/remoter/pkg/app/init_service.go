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

	execTimeout := time.Duration(p.fConfig.GetSettings().GetExecTimeoutMS()) * time.Millisecond
	p.fIncServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetIncoming(),
		Handler:     http.TimeoutHandler(mux, 2*execTimeout, "timeout"),
		ReadTimeout: (5 * time.Second),
	}
}
