package handler

import (
	"fmt"
	"image/png"
	"net/http"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/number571/go-peer/cmd/hlm/internal/app/state"
)

func QRPublicKeyPage(s state.IState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/qr/public_key" {
			NotFoundPage(s)(w, r)
			return
		}

		if !s.IsActive() {
			http.Redirect(w, r, "/sign/in", http.StatusFound)
			return
		}

		pubKey, err := s.GetClient().PubKey()
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
