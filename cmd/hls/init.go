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

	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/modules/client"
	"github.com/number571/go-peer/modules/crypto/asymmetric"
	"github.com/number571/go-peer/modules/filesystem"
	"github.com/number571/go-peer/modules/friends"
	"github.com/number571/go-peer/modules/inputter"
	"github.com/number571/go-peer/modules/logger"
	"github.com/number571/go-peer/modules/network"
	"github.com/number571/go-peer/modules/network/anonymity"
	"github.com/number571/go-peer/modules/network/conn_keeper"
	"github.com/number571/go-peer/modules/queue"
	"github.com/number571/go-peer/modules/storage"
	"github.com/number571/go-peer/modules/storage/database"

	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
)

func hlsDefaultInit() error {
	var (
		initOnly bool
		inputKey string
	)

	flag.BoolVar(&initOnly, "init", false, "run initialization only")
	flag.StringVar(&inputKey, "input-key", "", "[INSECURE] input private key from file")
	flag.Parse()

	gLogger = logger.NewLogger(os.Stdout, os.Stdout, os.Stdout)
	gConfig = config.NewConfig(hls_settings.CPathCFG)
	gEditor = config.NewEditor(&gConfig)

	var privKey asymmetric.IPrivKey
	switch inputKey {
	case "":
		privKey = initPrivKey(
			hls_settings.CPathSTG,
			[]byte(inputter.NewInputter("Password#Stg: ").Password()),
			[]byte(inputter.NewInputter("Password#Obj: ").Password()),
		)
		if privKey == nil {
			return fmt.Errorf("failed load private key")
		}
	default:
		privKeyStr, err := filesystem.OpenFile(inputKey).Read()
		if err != nil {
			return err
		}
		privKey = asymmetric.LoadRSAPrivKey(string(privKeyStr))
		if privKey == nil {
			return fmt.Errorf("private key is invalid")
		}
	}

	gLogger.Info(privKey.PubKey().String())
	if initOnly {
		os.Exit(0)
	}

	gNode = initNode(gConfig, gLogger, privKey)
	if err := gNode.Run(); err != nil {
		return err
	}

	gServerHTTP = initServerHTTP(gConfig, gLogger)
	gConnKeeper = conn_keeper.NewConnKeeper(
		conn_keeper.NewSettings(&conn_keeper.SSettings{
			FConnections: func() []string { return gConfig.Connections() },
			FDuration:    time.Minute,
		}),
		gNode.Network(),
	)
	if err := gConnKeeper.Run(); err != nil {
		return err
	}

	return nil
}

func initServerHTTP(cfg config.IConfig, logger logger.ILogger) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleIndexHTTP)
	mux.HandleFunc(hls_settings.CHandleConnects, handleConnectsHTTP)
	mux.HandleFunc(hls_settings.CHandleFriends, handleFriendsHTTP)
	mux.HandleFunc(hls_settings.CHandleOnline, handleOnlineHTTP)
	mux.HandleFunc(hls_settings.CHandlePubKey, handlePubKeyHTTP)
	mux.HandleFunc(hls_settings.CHandleRequest, handleRequestHTTP)

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
		// if node in client mode
		// then run endless loop
		if gConfig.Address().TCP() == "" {
			select {}
		}

		// run node in server mode
		gLogger.Info(fmt.Sprintf("TCP is listening [%s]...", gConfig.Address().TCP()))
		err := gNode.Network().Listen(gConfig.Address().TCP())
		if err != nil && !errors.Is(err, net.ErrClosed) {
			gLogger.Warning(err.Error())
		}
	}()

	return node
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
	answ := inputter.NewInputter("Private key by password not exist.\nGenerate new? [y/n]: ").String()
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
