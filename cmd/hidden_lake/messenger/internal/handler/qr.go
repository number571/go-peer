package handler

import (
	"fmt"
	"image/png"
	"net/http"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/app/state"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"
)

func QRPublicKeyPage(pStateManager state.IStateManager, pLogger logger.ILogger) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		if pR.URL.Path != "/qr/public_key" {
			NotFoundPage(pStateManager, pLogger)(pW, pR)
			return
		}

		if !pStateManager.StateIsActive() {
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogRedirect))
			http.Redirect(pW, pR, "/sign/in", http.StatusFound)
			return
		}

		myPubKey, _, err := pStateManager.GetClient().GetPubKey()
		if err != nil || !pStateManager.IsMyPubKey(myPubKey) {
			pLogger.PushWarn(logBuilder.WithMessage("get_public_key"))
			fmt.Fprint(pW, "error: read public key")
			return
		}

		qrCode, err := qr.Encode(myPubKey.ToString(), qr.L, qr.Auto)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("qr_encode"))
			fmt.Fprint(pW, "error: qrcode generate")
			return
		}

		qrCode, err = barcode.Scale(qrCode, 1024, 1024)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("qr_scale"))
			fmt.Fprint(pW, "error: qrcode scale")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		png.Encode(pW, qrCode)
	}
}
