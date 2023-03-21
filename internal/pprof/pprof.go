package pprof

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/logger"

	_ "net/http/pprof"
)

func RunPprofService(pService string) {
	logger := logger.NewLogger(logger.NewSettings(&logger.SSettings{
		FInfo: os.Stdout,
		FWarn: os.Stdout,
		FErro: os.Stderr,
	}))
	go runPprofService(logger, pService)
	time.Sleep(cWaitTime)
}

func runPprofService(pLogger logger.ILogger, pService string) {
	for i := 0; i < cRetriesNum; i++ {
		port := random.NewStdPRNG().GetUint64()%(4<<10) + 60000 // [60000;64096]
		pLogger.PushInfo(fmt.Sprintf("%s => pprof running on %d port", pService, port))

		err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			pLogger.PushWarn(fmt.Sprintf("%s => pprof service is closed", pService))
			return
		}
	}

	pLogger.PushErro(fmt.Sprintf("%s => pprof failed running", pService))
}
