package handler

import (
	"net/http"

	"github.com/number571/go-peer/cmd/hidden_lake/_template/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/_template/pkg/settings"
	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	testutils "github.com/number571/go-peer/test/_data"
)

func testRunService(addr string) *http.Server {
	mux := http.NewServeMux()

	// TODO: need implementation
	_ = &config.SConfig{}

	logger := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	mux.HandleFunc(settings.CHandleIndexPath, HandleIndexAPI(logger))

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}

func testNetworkMessageSettings() net_message.ISettings {
	return net_message.NewSettings(&net_message.SSettings{
		FNetworkKey:   testutils.TCNetworkKey,
		FWorkSizeBits: testutils.TCWorkSize,
	})
}
