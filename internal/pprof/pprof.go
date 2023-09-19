package pprof

import (
	"net/http"

	"net/http/pprof"
)

func InitPprofService(pAddr string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	server := &http.Server{
		Addr:    pAddr,
		Handler: mux,
	}

	return server
}

// func runPprofService(pLogger logger.ILogger, pService, pAddr string) {
// 	pLogger.PushInfo(fmt.Sprintf("%s => pprof running on %s", pService, pAddr))
// 	if err := http.ListenAndServe(pAddr, nil); err != nil {
// 		if errors.Is(err, http.ErrServerClosed) {
// 			pLogger.PushWarn(fmt.Sprintf("%s => pprof service is closed", pService))
// 			return
// 		}
// 		pLogger.PushErro(fmt.Sprintf("%s => pprof failed running", pService))
// 	}
// }
