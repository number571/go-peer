package handler

import (
	"crypto/sha256"
	"io"
	"io/fs"
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
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		page, err := strconv.Atoi(pR.URL.Query().Get("page"))
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("incorrect_page"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: incorrect page")
			return
		}

		result, err := getListFileInfo(pCfg, pStgPath, uint64(page))
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("open storage"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: open storage")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, result)
	}
}

func getListFileInfo(pCfg config.IConfig, pStgPath string, pPage uint64) ([]hlf_settings.SFileInfo, error) {
	pageOffset := pCfg.GetSettings().GetPageOffset()
	fileReader := pageOffset

	entries, err := os.ReadDir(pStgPath)
	if err != nil {
		return nil, err
	}

	files := make([]fs.DirEntry, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		files = append(files, e)
	}

	result := make([]hlf_settings.SFileInfo, 0, pageOffset)
	for i := (pPage * pageOffset); i < uint64(len(files)); i++ {
		if fileReader == 0 {
			break
		}
		fileReader--

		fileName := files[i].Name()
		fullPath := filepath.Join(pStgPath, fileName)

		result = append(result, hlf_settings.SFileInfo{
			FName: fileName,
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
