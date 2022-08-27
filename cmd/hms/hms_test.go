package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/client"
	"github.com/number571/go-peer/cmd/hms/config"
	"github.com/number571/go-peer/cmd/hms/database"
	"github.com/number571/go-peer/cmd/hms/hmc"
	hms_settings "github.com/number571/go-peer/cmd/hms/settings"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/payload"
	"github.com/number571/go-peer/utils/testutils"
)

const (
	tcN = 3
)

const (
	tcPathDB     = "hms_test.db"
	tcPathConfig = "hms_test.cfg"
)

var (
	tgHost = fmt.Sprintf("http://%s", testutils.TgAddrs[6])
)

func testHmsDefaultInit(dbPath, configPath string) {
	os.RemoveAll(dbPath)

	gDB = database.NewKeyValueDB(dbPath)
	gConfig = config.NewConfig(configPath)
}

func TestHMS(t *testing.T) {
	testHmsDefaultInit(tcPathDB, tcPathConfig)
	defer func() {
		gDB.Close()
		os.RemoveAll(tcPathDB)
		os.Remove(tcPathConfig)
	}()

	// server
	srv := testStartServerHTTP(t)
	defer srv.Close()

	// client push
	time.Sleep(200 * time.Millisecond)
	err := testClientDoPush()
	if err != nil {
		t.Error(err)
		return
	}

	// client size
	err = testClientDoSize()
	if err != nil {
		t.Error(err)
		return
	}

	// client load
	err = testClientDoLoad()
	if err != nil {
		t.Error(err)
		return
	}
}

func testStartServerHTTP(t *testing.T) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexPage)
	mux.HandleFunc("/size", sizePage)
	mux.HandleFunc("/load", loadPage)
	mux.HandleFunc("/push", pushPage)

	srv := &http.Server{
		Addr:    testutils.TgAddrs[6],
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}

func testClientDoPush() error {
	client := testNewClient()

	for i := 0; i < tcN; i++ {
		err := hmc.NewClient(
			hmc.NewBuilder(client),
			hmc.NewRequester(tgHost),
		).Push(
			client.PubKey(),
			payload.NewPayload(
				0x01,
				[]byte(fmt.Sprintf(testutils.TcBodyTemplate, i)),
			),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func testClientDoSize() error {
	client := testNewClient()

	size, err := hmc.NewClient(
		hmc.NewBuilder(client),
		hmc.NewRequester(tgHost),
	).Size()
	if err != nil {
		return err
	}

	if size != tcN {
		return fmt.Errorf("num(%d) != tcN(%d)", size, tcN)
	}

	return nil
}

func testClientDoLoad() error {
	client := testNewClient()

	for i := 0; i < tcN; i++ {
		msg, err := hmc.NewClient(
			hmc.NewBuilder(client),
			hmc.NewRequester(tgHost),
		).Load(uint64(i))
		if err != nil {
			return err
		}

		pubKey, pld, err := client.Decrypt(msg)
		if err != nil {
			panic(err)
		}

		if string(pld.Body()) != fmt.Sprintf(testutils.TcBodyTemplate, i) {
			return fmt.Errorf("body is not equal")
		}

		if pubKey.Address().String() != client.PubKey().Address().String() {
			return fmt.Errorf("public key is not equal")
		}
	}

	return nil
}

func testNewClient() client.IClient {
	return client.NewClient(
		client.NewSettings(&client.SSettings{
			FWorkSize:    hms_settings.CSizeWork,
			FMessageSize: hms_settings.CSizePack,
		}),
		asymmetric.LoadRSAPrivKey(testutils.TcPrivKey),
	)
}
