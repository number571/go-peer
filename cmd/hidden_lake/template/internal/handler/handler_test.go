package handler

import (
	"net/http"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/template/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/template/pkg/settings"
	"github.com/number571/go-peer/pkg/logger"
)

func testRunService(addr string) *http.Server {
	mux := http.NewServeMux()

	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FValue: "TODO",
		},
	}

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	mux.HandleFunc(settings.CHandleIndexPath, HandleIndexAPI(logger))
	mux.HandleFunc(settings.CHandleConfigSettingsPath, HandleConfigSettingsAPI(cfg, logger))

	srv := &http.Server{
		Addr:        addr,
		ReadTimeout: time.Second,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
	}

	go func() { _ = srv.ListenAndServe() }()
	return srv
}
