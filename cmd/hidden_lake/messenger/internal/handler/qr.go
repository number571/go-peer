package handler

import (
	"fmt"
	"image/png"
	"net/http"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/app/state"
)

func QRPublicKeyPage(pState state.IState) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.URL.Path != "/qr/public_key" {
			NotFoundPage(pState)(pW, pR)
			return
		}

		if !pState.IsActive() {
			http.Redirect(pW, pR, "/sign/in", http.StatusFound)
			return
		}

		pubKey, err := pState.GetClient().Service().GetPubKey()
		if err != nil {
			fmt.Fprint(pW, "error: read public key")
			return
		}

		qrCode, err := qr.Encode(pubKey.ToString(), qr.L, qr.Auto)
		if err != nil {
			fmt.Fprint(pW, "error: qrcode generate")
			return
		}

		qrCode, err = barcode.Scale(qrCode, 1024, 1024)
		if err != nil {
			fmt.Fprint(pW, "error: qrcode scale")
			return
		}

		png.Encode(pW, qrCode)
	}
}
