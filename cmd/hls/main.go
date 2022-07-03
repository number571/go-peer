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
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/crypto/random"
	"github.com/number571/go-peer/local/selector"
	"github.com/number571/go-peer/netanon"
	"github.com/number571/go-peer/settings"
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

	// set response route
	gNode.WithRouter(func() []asymmetric.IPubKey {
		randSizeRoute := random.NewStdPRNG().Uint64() % hls_settings.CSizeRoute
		return selector.NewSelector(nodesInOnline(gNode)).
			Shuffle().
			Return(randSizeRoute)
	})

	// turn on pseudo packages
	gNode.Pseudo().Switch(true)

	// turn on f2f mode
	gNode.F2F().Switch(gConfig.F2F().Status())
	for _, pubKey := range gConfig.F2F().PubKeys() {
		gNode.F2F().Append(pubKey)
	}

	// turn on online checker
	isOnline := gConfig.OnlineChecker().Status()
	gNode.Online().Switch(isOnline)
	gNode.Checker().Switch(isOnline)
	for _, pubKey := range gConfig.OnlineChecker().PubKeys() {
		gNode.Checker().Append(pubKey)
	}

	// connect to open nodes
	for _, address := range gConfig.Connections() {
		if address == gConfig.Address().HLS() {
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
		sett := gNode.Client().Settings()
		tchk := time.Duration(sett.Get(settings.CTimePing))

		for {
			time.Sleep(time.Second * tchk)
			for _, address := range gConfig.Connections() {
				if address == gConfig.Address().HLS() {
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
	mux.HandleFunc("/status", pageStatus)
	mux.HandleFunc("/message", pageMessage)

	srv := &http.Server{
		Addr:    gConfig.Address().HTTP(),
		Handler: mux,
	}

	go func() {
		gLogger.Info(fmt.Sprintf("HTTP is listening [%s]...", gConfig.Address().HTTP()))
		err := srv.ListenAndServe()
		if err != nil {
			gLogger.Warning(err.Error())
		}
	}()

	go func() {
		gLogger.Info(fmt.Sprintf("HLS is listening [%s]...", gConfig.Address().HLS()))

		// if node in client mode
		// then run endless loop
		if gConfig.Address().HLS() == "" {
			select {}
		}

		// run node in server mode
		err = gNode.Network().Listen(gConfig.Address().HLS())
		if err != nil {
			gLogger.Warning(err.Error())
		}
	}()

	<-shutdown
	fmt.Println("Shutting down...")
	utils.CloseAll([]utils.ICloser{srv, gDB, gNode})
}

func nodesInOnline(node netanon.INode) []asymmetric.IPubKey {
	inOnline := []asymmetric.IPubKey{}
	for _, info := range gNode.Checker().ListWithInfo() {
		if !info.Online() {
			continue
		}
		inOnline = append(inOnline, info.PubKey())
	}
	return inOnline
}
