package handler

import (
	"fmt"
	"image/png"
	"net/http"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	hls_client "github.com/number571/go-peer/cmd/hls/pkg/client"
)

func QRPublicKeyPage(client hls_client.IClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/qr/public_key" {
			NotFoundPage(w, r)
			return
		}

		pubKey, err := client.PubKey()
		if err != nil {
			fmt.Fprint(w, "error: read public key")
			return
		}

		qrCode, err := qr.Encode(pubKey.String(), qr.L, qr.Auto)
		if err != nil {
			fmt.Fprint(w, "error: qrcode generate")
			return
		}

		qrCode, err = barcode.Scale(qrCode, 1024, 1024)
		if err != nil {
			fmt.Fprint(w, "error: qrcode scale")
			return
		}

		png.Encode(w, qrCode)
	}
}
