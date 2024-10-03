package handler

import (
	"context"
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

	if _, err := client.EncryptMessage(context.Background(), asymmetric.LoadRSAPubKey(testutils.TgPubKeys[0]), payload.NewPayload64(1, []byte{123})); err == nil {
		t.Error("success encrypt message with unknown host")
		return
	}

	pld := payload.NewPayload32(testutils.TcHead, []byte(testutils.TcBody))
	sett := message.NewConstructSettings(&message.SConstructSettings{
		FSettings: message.NewSettings(&message.SSettings{
			FWorkSizeBits: testutils.TCWorkSize,
		}),
	})
	if _, _, err := client.DecryptMessage(context.Background(), message.NewMessage(sett, pld)); err == nil {
		t.Error("success decrypt message with unknown host")
		return
	}

	if _, err := client.GetIndex(context.Background()); err == nil {
		t.Error("success get index with unknown host")
		return
	}

	if _, err := client.GetSettings(context.Background()); err == nil {
		t.Error("success get settings with unknown host")
		return
	}

	if _, err := client.GetPubKey(context.Background()); err == nil {
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

	title, err := hleClient.GetIndex(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if title != settings.CServiceFullName {
		t.Error("incorrect title pattern")
		return
	}
}
