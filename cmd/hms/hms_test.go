package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/number571/go-peer/cmd/hms/config"
	"github.com/number571/go-peer/cmd/hms/database"
	"github.com/number571/go-peer/cmd/hms/hmc"
	"github.com/number571/go-peer/crypto/asymmetric"
	"github.com/number571/go-peer/local/client"
	"github.com/number571/go-peer/settings/testutils"
)

const (
	tcPathDB     = "hms_test.db"
	tcPathConfig = "hms_test.cfg"
)

const (
	tcN = 3

	tcHost           = "http://" + tcServiceAddress
	tcServiceAddress = "localhost:8080"
	tcBodyOfMessage  = "hello, world! (%d)"

	tcPrivKeyReceiver = `Priv(go-peer/rsa){3082025E02010002818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C702030100010281810083866C4CA38EDAACE6B62A69A8C5682FD24F136A2E081C34F5AFB89372737AE3D052000317C7A2C9164180DD8E09E53C94F88341DFA8BD275E594CBAB9D4B008E1FE2D613D35202E841858BC665C0338221F34D9F143D60A5C2C4625459DAD3C0E3592F6B32D3E4105AB713CE42E73C44F10687954402E7A8D2952CC1C4589B9024100EB57FACD3A75AD2BA0BBF2152BF3760CFE45F78731C3B98D770DD790082753E4697CE8927112632F7BE86880121F4A08880DE7C45D16EB8D76E72214768D517B024100D8821E8FCF0DC72C5DD63A4A39CCC42601E1553022B75C9D01EF3DA2F706081E694C98E684BA5482B6E5975F6B371DE6E81AD42CB74A7CFA52A6D5522E0D4625024100BFDD211DF17C006AE206778CE520FDEC07DC98B9424BF3D92DE73E07316E86895FAAB29CB8CC29CA8B74E4C50C812FC516CE675602226E750D2BCFEFE8DABB4302407843AF1E4AF16855A8BA3B1EC8048A606262FCA30465BE3828BEF009FA158BA4F8F0E76E05044BB5604B204E8C8BCD3C5A69ACBA3A06526DEA4369F380493751024100EAE2ADA08E39EB52C8314DBB0F16A087DD9AE4BA3DADCD3BE515EA4193F83E62066ECDB3BE47CB377AE7F5480141FF60C20AAE818B3CDFAAA6244D97FB09FF0D}`
)

func testHmsDefaultInit(dbPath, configPath string) {
	os.RemoveAll(dbPath)

	gSettings = testutils.NewSettings()
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
		Addr:    tcServiceAddress,
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}

func testClientDoPush() error {
	priv := asymmetric.LoadRSAPrivKey(tcPrivKeyReceiver)
	client := client.NewClient(priv, gSettings)

	for i := 0; i < tcN; i++ {
		err := hmc.NewClient(
			hmc.NewBuiler(client),
			hmc.NewRequester(tcHost),
		).Push(
			client.PubKey(),
			[]byte(fmt.Sprintf(tcBodyOfMessage, i)),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func testClientDoSize() error {
	priv := asymmetric.LoadRSAPrivKey(tcPrivKeyReceiver)
	client := client.NewClient(priv, gSettings)

	size, err := hmc.NewClient(
		hmc.NewBuiler(client),
		hmc.NewRequester(tcHost),
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
	priv := asymmetric.LoadRSAPrivKey(tcPrivKeyReceiver)
	client := client.NewClient(priv, gSettings)

	for i := 0; i < tcN; i++ {
		msg, err := hmc.NewClient(
			hmc.NewBuiler(client),
			hmc.NewRequester(tcHost),
		).Load(uint64(i))
		if err != nil {
			return err
		}

		body := msg.Body().Data()
		if string(body) != fmt.Sprintf(tcBodyOfMessage, i) {
			return fmt.Errorf("body is not equal")
		}

		pubKey := asymmetric.LoadRSAPubKey(msg.Head().Sender())
		if pubKey.Address().String() != client.PubKey().Address().String() {
			return fmt.Errorf("public key is not equal")
		}
	}

	return nil
}
