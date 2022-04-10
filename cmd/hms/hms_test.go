package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/number571/go-peer/cmd/hms/database"
	"github.com/number571/go-peer/cmd/hms/utils"
	"github.com/number571/go-peer/crypto"
	"github.com/number571/go-peer/encoding"
	"github.com/number571/go-peer/local"
	"github.com/number571/go-peer/settings/testutils"
)

var (
	tgPubKeyReceiver = crypto.LoadPubKey(tcPubKeyReceiver)
)

const (
	tcPathDB = "hms_test.db"
)

const (
	tcN        = 3
	tcAKeySize = 1024

	tcServiceAddress  = "localhost:8081"
	tcPatternTitleHMS = "store-message"
	tcBodyOfMessage   = "hello, world!"

	tcPrivKeyReceiver = `Priv(go-peer\rsa){3082025E02010002818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C702030100010281810083866C4CA38EDAACE6B62A69A8C5682FD24F136A2E081C34F5AFB89372737AE3D052000317C7A2C9164180DD8E09E53C94F88341DFA8BD275E594CBAB9D4B008E1FE2D613D35202E841858BC665C0338221F34D9F143D60A5C2C4625459DAD3C0E3592F6B32D3E4105AB713CE42E73C44F10687954402E7A8D2952CC1C4589B9024100EB57FACD3A75AD2BA0BBF2152BF3760CFE45F78731C3B98D770DD790082753E4697CE8927112632F7BE86880121F4A08880DE7C45D16EB8D76E72214768D517B024100D8821E8FCF0DC72C5DD63A4A39CCC42601E1553022B75C9D01EF3DA2F706081E694C98E684BA5482B6E5975F6B371DE6E81AD42CB74A7CFA52A6D5522E0D4625024100BFDD211DF17C006AE206778CE520FDEC07DC98B9424BF3D92DE73E07316E86895FAAB29CB8CC29CA8B74E4C50C812FC516CE675602226E750D2BCFEFE8DABB4302407843AF1E4AF16855A8BA3B1EC8048A606262FCA30465BE3828BEF009FA158BA4F8F0E76E05044BB5604B204E8C8BCD3C5A69ACBA3A06526DEA4369F380493751024100EAE2ADA08E39EB52C8314DBB0F16A087DD9AE4BA3DADCD3BE515EA4193F83E62066ECDB3BE47CB377AE7F5480141FF60C20AAE818B3CDFAAA6244D97FB09FF0D}`
	tcPubKeyReceiver  = `Pub(go-peer\rsa){30818902818100C709DA63096CEDBA0DD6B5DD9465B412268C00509757A8EBD9096E17BEEC17C25A3A8F246E1591554CD214F4B27254EFA811F8BE441A03B37B3C8B390484C74C2294A4C895AA925D723E0065A877D4502CC010996863821E7348348E4E96CDD4CB7A852B2E2853C8FDEE556C4F89F6C3295EAC00DAEE86DD94E25F9703F368C70203010001}`
)

func testHmsDefaultInit(path string) {
	os.RemoveAll(tcPathDB)

	gSettings = testutils.NewSettings()
	gDB = database.NewKeyValueDB(path)
}

func TestHMS(t *testing.T) {
	testHmsDefaultInit(tcPathDB)
	defer os.RemoveAll(tcPathDB)

	// server
	srv := testStartServerHTTP(t)
	defer srv.Close()

	// client push
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
	priv := crypto.NewPrivKey(tcAKeySize)
	client := local.NewClient(priv, gSettings)

	hashRecv := crypto.NewHasher(tgPubKeyReceiver.Bytes()).Bytes()

	for i := 0; i < tcN; i++ {
		encMsg, _ := client.Encrypt(
			local.NewRoute(tgPubKeyReceiver),
			local.NewMessage(
				[]byte(tcPatternTitleHMS),
				[]byte(tcBodyOfMessage),
			),
		)

		request := struct {
			Receiver []byte `json:"receiver"`
			Package  []byte `json:"package"`
		}{
			Receiver: hashRecv,
			Package:  encMsg.ToPackage().Bytes(),
		}

		resp, err := http.Post(
			"http://"+tcServiceAddress+"/push",
			"application/json",
			bytes.NewReader(utils.Serialize(request)),
		)
		if err != nil {
			return err
		}

		var response struct {
			Result []byte `json:"result"`
			Return int    `json:"return"`
		}
		json.NewDecoder(resp.Body).Decode(&response)

		if response.Return != cErrorNone {
			return fmt.Errorf("%v", string(response.Result))
		}
	}

	return nil
}

func testClientDoSize() error {
	hashRecv := crypto.NewHasher(tgPubKeyReceiver.Bytes()).Bytes()

	request := struct {
		Receiver []byte `json:"receiver"`
	}{
		Receiver: hashRecv,
	}

	resp, err := http.Post(
		"http://"+tcServiceAddress+"/size",
		"application/json",
		bytes.NewReader(utils.Serialize(request)),
	)
	if err != nil {
		return err
	}

	var response struct {
		Result []byte `json:"result"`
		Return int    `json:"return"`
	}
	json.NewDecoder(resp.Body).Decode(&response)

	if response.Return != cErrorNone {
		return fmt.Errorf("%v", string(response.Result))
	}

	num := encoding.BytesToUint64(response.Result)
	if num != tcN {
		return fmt.Errorf("num(%d) != tcN(%d)", num, tcN)
	}

	return nil
}

func testClientDoLoad() error {
	priv := crypto.LoadPrivKey(tcPrivKeyReceiver)
	client := local.NewClient(priv, gSettings)

	hashRecv := crypto.NewHasher(tgPubKeyReceiver.Bytes()).Bytes()

	for i := 0; i < tcN; i++ {
		request := struct {
			Receiver []byte `json:"receiver"`
			Index    uint64
		}{
			Receiver: hashRecv,
			Index:    uint64(i),
		}

		resp, err := http.Post(
			"http://"+tcServiceAddress+"/load",
			"application/json",
			bytes.NewReader(utils.Serialize(request)),
		)
		if err != nil {
			return err
		}

		var response struct {
			Result []byte `json:"result"`
			Return int    `json:"return"`
		}
		json.NewDecoder(resp.Body).Decode(&response)

		if response.Return != cErrorNone {
			return fmt.Errorf("%v", string(response.Result))
		}

		msg := local.LoadPackage(response.Result).ToMessage()
		if msg == nil {
			return fmt.Errorf("message is nil")
		}

		msg, title := client.Decrypt(msg)
		if string(title) != tcPatternTitleHMS {
			return fmt.Errorf("title is not equal")
		}

		if string(msg.Body().Data()) != tcBodyOfMessage {
			return fmt.Errorf("message is not equal")
		}
	}

	return nil
}
