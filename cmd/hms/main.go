// HMS - Hidden Message Service
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
)

func main() {
	err := hmsDefaultInit()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexPage)
	mux.HandleFunc(hms_settings.CSizePath, sizePage)
	mux.HandleFunc(hms_settings.CLoadPath, loadPage)
	mux.HandleFunc(hms_settings.CPushSize, pushPage)

	srv := &http.Server{
		Addr:    gConfig.Address(),
		Handler: mux,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	<-shutdown
	fmt.Println("Shutting down...")

	srv.Close()
	gDB.Close()
}
