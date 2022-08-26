// HLS - Hidden Lake Service
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
	"github.com/number571/go-peer/utils"
)

func main() {
	// read config, database, key
	err := hlsDefaultInit()
	if err != nil {
		gLogger.Error(err.Error())
		os.Exit(1)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// set handle functions
	gNode.Handle(hls_settings.CHeaderHLS, routeHLS)

	// turn on f2f mode
	for _, pubKey := range gConfig.Friends() {
		gNode.F2F().Append(pubKey)
	}

	// connect to open nodes
	for _, address := range gConfig.Connections() {
		if address == gConfig.Address().TCP() {
			gLogger.Warning(fmt.Sprintf("used own address '%s'", address))
			continue
		}
		conn := gNode.Network().Connect(address)
		if conn == nil {
			gLogger.Warning("conn is nil")
			continue
		}
		gLogger.Info(fmt.Sprintf("connected to '%s'", address))
	}

	// network checker
	go func() {
		for {
			time.Sleep(time.Minute)
			for _, address := range gConfig.Connections() {
				if address == gConfig.Address().TCP() {
					continue
				}
				for _, conn := range gNode.Network().Connections() {
					if conn.Socket().LocalAddr().String() == address {
						continue
					}
				}
				conn := gNode.Network().Connect(address)
				if conn == nil {
					gLogger.Warning(fmt.Sprintf("conn is nil %s", address))
					continue
				}
				gLogger.Info(fmt.Sprintf("connected to '%s'", address))
			}
		}
	}()

	// HTTP client
	mux := http.NewServeMux()

	mux.HandleFunc("/", pageIndex)
	mux.HandleFunc("/send", pageSend)

	srv := &http.Server{
		Addr:    gConfig.Address().HTTP(),
		Handler: mux,
	}

	go func() {
		gLogger.Info(fmt.Sprintf("HTTP is listening [%s]...", gConfig.Address().HTTP()))
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			gLogger.Warning(err.Error())
		}
	}()

	go func() {
		gLogger.Info(fmt.Sprintf("TCP is listening [%s]...", gConfig.Address().TCP()))

		// if node in client mode
		// then run endless loop
		if gConfig.Address().TCP() == "" {
			select {}
		}

		// run node in server mode
		err = gNode.Network().Listen(gConfig.Address().TCP())
		if err != nil {
			gLogger.Warning(err.Error())
		}
	}()

	<-shutdown
	fmt.Println("Shutting down...")
	utils.CloseAll([]utils.ICloser{srv, gNode})
}
