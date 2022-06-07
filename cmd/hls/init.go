package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/number571/go-peer/cmd/hls/config"
	"github.com/number571/go-peer/cmd/hls/database"
	"github.com/number571/go-peer/cmd/hls/logger"
	"github.com/number571/go-peer/cmd/hls/utils"
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/local"
	"github.com/number571/go-peer/network"
	"github.com/number571/go-peer/settings"
	"github.com/number571/go-peer/storage"

	hls_settings "github.com/number571/go-peer/cmd/hls/settings"
)

func hlsDefaultInit() error {
	var (
		initOnly bool
	)

	flag.BoolVar(&initOnly, "init", false, "run initialization only")
	flag.Parse()

	gPPrivKey = crypto.NewPrivKey(hls_settings.CAKeySize)
	gLogger = logger.NewLogger(os.Stdout, os.Stdout, os.Stdout)
	gConfig = config.NewConfig("hls.cfg")
	gDB = database.NewKeyValueDB("hls.db")

	sett := hls_settings.NewSettings()
	privKey := getPrivKey(
		sett,
		"hls.stg",
		[]byte(utils.InputString("Storage password: ")),
		[]byte(utils.InputString("Object password: ")),
	)
	if privKey == nil {
		return fmt.Errorf("failed load private key")
	}

	gNode = network.NewNode(local.NewClient(privKey, sett))
	if gNode == nil {
		return fmt.Errorf("failed create client node")
	}

	if initOnly {
		os.Exit(0)
	}

	return nil
}

func getPrivKey(sett settings.ISettings, filepath string, storageKey, objectKey []byte) crypto.IPrivKey {
	storage := storage.NewCryptoStorage(
		sett,
		filepath,
		storageKey,
	)
	if storage == nil {
		return nil
	}

	bpriv, err := storage.Get(objectKey)
	if err == nil {
		return crypto.LoadPrivKey(bpriv)
	}

	priv := crypto.NewPrivKey(hls_settings.CAKeySize)
	err = storage.Set(objectKey, priv.Bytes())
	if err != nil {
		return nil
	}

	return priv
}
