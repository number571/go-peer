package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/number571/go-peer/client"
	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/cmd/hls/database"
	"github.com/number571/go-peer/cmd/hls/logger"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/friends"
	"github.com/number571/go-peer/netanon"
	"github.com/number571/go-peer/network"
	"github.com/number571/go-peer/queue"
	"github.com/number571/go-peer/storage"
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
	gConfig = config.NewConfig("hls.cfg")
	gDB = database.NewKeyValueDB("hls.db")

	privKey := getPrivKey(
		"hls.stg",
		[]byte(utils.NewInput("Password#Stg: ").Password()),
		[]byte(utils.NewInput("Password#Obj: ").Password()),
	)
	if privKey == nil {
		return fmt.Errorf("failed load private key")
	}

	client := client.NewClient(
		client.NewSettings(hls_settings.CWorkSize, hls_settings.CRandBytes),
		privKey,
	)
	gNode = netanon.NewNode(
		netanon.NewSettings(
			1,
			2,
			hls_settings.CWaitTime*time.Second,
		),
		client,
		network.NewNode(network.NewSettings(
			hls_settings.CPackSize,
			10,   // retryNum for get message
			1024, // capacity for hash storage
			hls_settings.CMaxConns,
			hls_settings.CMaxMsgs,
			5*time.Second, // timeWait for request
		)),
		queue.NewQueue(
			queue.NewSettings(
				hls_settings.CQueueSize,
				hls_settings.CQueuePull,
				hls_settings.CPackSize,
				hls_settings.CQueueTime*time.Millisecond,
			),
			client,
		),
		friends.NewF2F(),
		func() []asymmetric.IPubKey {
			return nil // TODO
		},
	)

	gLogger.Info(privKey.PubKey().String())
	if initOnly {
		os.Exit(0)
	}

	return nil
}

func getPrivKey(filepath string, storageKey, objectKey []byte) asymmetric.IPrivKey {
	// create/open storage
	storage := storage.NewCryptoStorage(
		filepath,
		storageKey,
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
