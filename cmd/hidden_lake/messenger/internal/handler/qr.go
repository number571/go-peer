package handler

import (
	"fmt"
	"image/png"
	"net/http"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/number571/go-peer/cmd/hidden_lake/messenger/internal/config"
	hlm_settings "github.com/number571/go-peer/cmd/hidden_lake/messenger/pkg/settings"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"
)

func QRPublicKeyPage(pLogger logger.ILogger, pCfg config.IConfig) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		if pR.URL.Path != "/qr/public_key" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		myPubKey, err := getClient(pCfg).GetPubKey()
		if err != nil {
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
