package handler

import (
	"net/http"
	"testing"
	"time"

	hle_client "github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/client"
	"github.com/number571/go-peer/cmd/hidden_lake/helpers/encryptor/pkg/settings"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/go-peer/test/utils"
)

func TestErrorsAPI(t *testing.T) {
	t.Parallel()

	client := hle_client.NewClient(
		hle_client.NewRequester(
			"http://"+testutils.TcUnknownHost,
			&http.Client{Timeout: time.Second},
			testNetworkMessageSettings(),
		),
	)

	if _, err := client.EncryptMessage(asymmetric.LoadRSAPubKey(testutils.TgPubKeys[0]), []byte{123}); err == nil {
		t.Error("success encrypt message with unknown host")
		return
	}

	pld := payload.NewPayload(uint64(testutils.TcHead), []byte(testutils.TcBody))
	sett := message.NewSettings(&message.SSettings{
		FWorkSizeBits: testutils.TCWorkSize,
	})
	if _, _, err := client.DecryptMessage(message.NewMessage(sett, pld, 1)); err == nil {
		t.Error("success decrypt message with unknown host")
		return
	}

	if _, err := client.GetIndex(); err == nil {
		t.Error("success get index with unknown host")
		return
	}

	if _, err := client.GetSettings(); err == nil {
		t.Error("success get settings with unknown host")
		return
	}

	if _, err := client.GetPubKey(); err == nil {
		t.Error("success get pub key with unknown host")
		return
	}
}

func TestHandleIndexAPI(t *testing.T) {
	t.Parallel()

	service := testRunService(testutils.TgAddrs[54])
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hleClient := hle_client.NewClient(
		hle_client.NewRequester(
			"http://"+testutils.TgAddrs[54],
			&http.Client{Timeout: time.Second / 2},
			testNetworkMessageSettings(),
		),
	)

	title, err := hleClient.GetIndex()
	if err != nil {
		t.Error(err)
		return
	}

	if title != settings.CServiceFullName {
		t.Error("incorrect title pattern")
		return
	}
}
