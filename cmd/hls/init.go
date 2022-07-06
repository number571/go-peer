package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/cmd/hls/database"
	"github.com/number571/go-peer/cmd/hls/logger"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/local/client"
	"github.com/number571/go-peer/netanon"
	"github.com/number571/go-peer/settings"
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

	sett := hls_settings.NewSettings()
	privKey := getPrivKey(
		sett,
		"hls.stg",
		[]byte(utils.NewInput(sett, "Password#Stg: ").Password()),
		[]byte(utils.NewInput(sett, "Password#Obj: ").Password()),
	)
	if privKey == nil {
		return fmt.Errorf("failed load private key")
	}

	gNode = netanon.NewNode(client.NewClient(sett, privKey))
	if gNode == nil {
		return fmt.Errorf("failed create client node")
	}

	if initOnly {
		os.Exit(0)
	}

	return nil
}

func getPrivKey(sett settings.ISettings, filepath string, storageKey, objectKey []byte) asymmetric.IPrivKey {
	fileAlreadyExist := utils.OpenFile(filepath).IsExist()

	storage := storage.NewCryptoStorage(
		sett,
		filepath,
		storageKey,
	)
	if storage == nil {
		return nil
	}

	if fileAlreadyExist {
		bpriv, err := storage.Get(objectKey)
		if err != nil {
			return nil
		}
		return asymmetric.LoadRSAPrivKey(bpriv)
	}

	priv := asymmetric.NewRSAPrivKey(hls_settings.CAKeySize)
	err := storage.Set(objectKey, priv.Bytes())
	if err != nil {
		return nil
	}

	return priv
}
