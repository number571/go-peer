package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/number571/go-peer/client"
	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/cmd/hls/logger"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/database"
	"github.com/number571/go-peer/friends"
	"github.com/number571/go-peer/netanon"
	"github.com/number571/go-peer/network"
	"github.com/number571/go-peer/queue"
	"github.com/number571/go-peer/storage"
	"github.com/number571/go-peer/testutils"
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

	privKey := getPrivKey(
		"hls.stg",
		[]byte(utils.NewInput("Password#Stg: ").Password()),
		[]byte(utils.NewInput("Password#Obj: ").Password()),
	)
	if privKey == nil {
		return fmt.Errorf("failed load private key")
	}

	gNode = netanon.NewNode(
		netanon.NewSettings(&netanon.SSettings{
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
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey),
			),
		),
		friends.NewF2F(),
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
