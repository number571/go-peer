package handler

import (
	"crypto/sha256"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/config"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"

	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func HandleIncomigListHTTP(pLogger logger.ILogger, pCfg config.IConfig, pStgPath string) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hls_settings.CHeaderResponseMode, hls_settings.CHeaderResponseModeON)

		logBuilder := http_logger.NewLogBuilder(hlf_settings.CServiceName, pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		page, err := strconv.Atoi(pR.URL.Query().Get("page"))
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("incorrect_page"))
			api.Response(pW, http.StatusBadRequest, "failed: incorrect page")
			return
		}

		result, err := getListFileInfo(pCfg, pStgPath, uint64(page))
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("open storage"))
			api.Response(pW, http.StatusInternalServerError, "failed: open storage")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		api.Response(pW, http.StatusOK, result)
	}
}

func getListFileInfo(pCfg config.IConfig, pStgPath string, pPage uint64) ([]hlf_settings.SFileInfo, error) {
	pageOffset := pCfg.GetSettings().GetPageOffset()

	entries, err := os.ReadDir(pStgPath)
	if err != nil {
		return nil, err
	}
	lenEntries := uint64(len(entries))

	result := make([]hlf_settings.SFileInfo, 0, lenEntries)
	for i := (pPage * pageOffset); i < lenEntries; i++ {
		e := entries[i]
		if e.IsDir() {
			continue
		}
		if i != (pPage*pageOffset) && i%pageOffset == 0 {
			break
		}
		fullPath := filepath.Join(pStgPath, e.Name())
		result = append(result, hlf_settings.SFileInfo{
			FName: e.Name(),
			FHash: getFileHash(fullPath),
			FSize: getFileSize(fullPath),
		})
	}
	return result, nil
}

func getFileSize(filename string) uint64 {
	stat, _ := os.Stat(filename)
	return uint64(stat.Size())
}

func getFileHash(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}

	return encoding.HexEncode(h.Sum(nil))
}
