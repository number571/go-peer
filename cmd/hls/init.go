package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/number571/go-peer/client"
	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/friends"
	"github.com/number571/go-peer/logger"
	"github.com/number571/go-peer/network"
	"github.com/number571/go-peer/network/anonymity"
	"github.com/number571/go-peer/network/conn_keeper"
	"github.com/number571/go-peer/queue"
	"github.com/number571/go-peer/storage"
	"github.com/number571/go-peer/storage/database"
	"github.com/number571/go-peer/utils"

	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
)

func hlsDefaultInit() error {
	var (
		initOnly bool
	)

	flag.BoolVar(&initOnly, "init", false, "run initialization only")
	flag.Parse()

	gLogger = logger.NewLogger(os.Stdout, os.Stdout, os.Stdout)
	gConfig = config.NewConfig(hls_settings.CPathCFG)

	privKey := initPrivKey(
		hls_settings.CPathSTG,
		[]byte(utils.NewInput("Password#Stg: ").Password()),
		[]byte(utils.NewInput("Password#Obj: ").Password()),
	)
	if privKey == nil {
		return fmt.Errorf("failed load private key")
	}

	gLogger.Info(privKey.PubKey().String())
	if initOnly {
		os.Exit(0)
	}

	gServerHTTP = initServerHTTP(gConfig, gLogger)
	gNode = initNode(gConfig, gLogger, privKey)
	if err := gNode.Run(); err != nil {
		return err
	}

	return conn_keeper.NewConnKeeper(
		conn_keeper.NewSettings(
			&conn_keeper.SSettings{
				FConnections: gConfig.Connections(),
				FDuration:    time.Minute,
			},
		),
		gNode.Network(),
	).Run()
}

func initNode(cfg config.IConfig, logger logger.ILogger, privKey asymmetric.IPrivKey) anonymity.INode {
	node := anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FRetryEnqueue: hls_settings.CRetryEnqueue,
			FTimeWait:     hls_settings.CWaitTime,
		}),
		database.NewLevelDB(
			database.NewSettings(&database.SSettings{
				FPath:    hls_settings.CPathDB,
				FHashing: true,
			}),
		),
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FRetryNum:    hls_settings.CNetworkRetry,
				FCapacity:    hls_settings.CNetworkCapacity,
				FMessageSize: hls_settings.CMessageSize,
				FMaxConns:    hls_settings.CNetworkMaxConns,
				FMaxMessages: hls_settings.CNetworkMaxMessages,
				FTimeWait:    hls_settings.CNetworkWaitTime,
			}),
		),
		queue.NewQueue(
			queue.NewSettings(&queue.SSettings{
				FCapacity:     hls_settings.CQueueCapacity,
				FPullCapacity: hls_settings.CQueuePullCapacity,
				FDuration:     hls_settings.CQueueDuration,
			}),
			client.NewClient(
				client.NewSettings(&client.SSettings{
					FWorkSize:    hls_settings.CWorkSize,
					FMessageSize: hls_settings.CMessageSize,
				}),
				privKey,
			),
		),
		func() friends.IF2F {
			f2f := friends.NewF2F()
			for _, pubKey := range cfg.Friends() {
				f2f.Append(pubKey)
			}
			return f2f
		}(),
	).Handle(hls_settings.CHeaderHLS, handleTCP)

	go func() {
		gLogger.Info(fmt.Sprintf("TCP is listening [%s]...", gConfig.Address().TCP()))

		// if node in client mode
		// then run endless loop
		if gConfig.Address().TCP() == "" {
			select {}
		}

		// run node in server mode
		err := gNode.Network().Listen(gConfig.Address().TCP())
		if err != nil && !errors.Is(err, net.ErrClosed) {
			gLogger.Warning(err.Error())
		}
	}()

	return node
}

func initServerHTTP(cfg config.IConfig, logger logger.ILogger) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleIndexHTTP)
	mux.HandleFunc("/request", handleRequestHTTP)

	srv := &http.Server{
		Addr:    cfg.Address().HTTP(),
		Handler: mux,
	}

	go func() {
		logger.Info(fmt.Sprintf("HTTP is listening [%s]...", cfg.Address().HTTP()))
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Warning(err.Error())
		}
	}()

	return srv
}

func initPrivKey(filepath string, storageKey, objectKey []byte) asymmetric.IPrivKey {
	// create/open storage
	storage := storage.NewCryptoStorage(
		filepath,
		storageKey,
		hls_settings.CWorkSize,
	)
	if storage == nil {
		return nil
	}

	// get private key
	bpriv, err := storage.Get(objectKey)
	if err == nil {
		return asymmetric.LoadRSAPrivKey(bpriv)
	}

	// private key not exist
	answ := utils.NewInput("Private key by password not exist.\nGenerate new? [y/n]: ").String()
	switch strings.ToLower(answ) {
	case "y", "yes":
		// generate private key
		priv := asymmetric.NewRSAPrivKey(hls_settings.CAKeySize)
		err := storage.Set(objectKey, priv.Bytes())
		if err != nil {
			panic(err)
		}
		return priv
	case "n", "no":
		// exit from program
		return nil
	default:
		// undefined answer
		panic("input answer not equal [y/yes, n/no]")
	}
}
