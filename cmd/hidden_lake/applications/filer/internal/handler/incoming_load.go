package handler

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filer/internal/config"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"

	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filer/pkg/settings"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func HandleIncomigLoadHTTP(pLogger logger.ILogger, pCfg config.IConfig, pPathTo string) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hls_settings.CHeaderResponseMode, hls_settings.CHeaderResponseModeON)

		logBuilder := http_logger.NewLogBuilder(hlf_settings.CServiceName, pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		query := pR.URL.Query()

		name := filepath.Base(query.Get("name"))
		if name != query.Get("name") {
			pLogger.PushWarn(logBuilder.WithMessage("got_another_name"))
			api.Response(pW, http.StatusConflict, "failed: got another name")
			return
		}

		chunk, err := strconv.Atoi(query.Get("chunk"))
		if err != nil || chunk < 0 {
			pLogger.PushWarn(logBuilder.WithMessage("incorrect_chunk"))
			api.Response(pW, http.StatusBadRequest, "failed: incorrect chunk")
			return
		}

		fPubKey := asymmetric.LoadRSAPubKey(pR.Header.Get(hls_settings.CHeaderPublicKey))
		if fPubKey == nil {
			panic("public key is nil (invalid data from HLS)!")
		}

		fullPath := fmt.Sprintf("%s/%s/%s", pPathTo, hlf_settings.CPathSTG, name)
		stat, err := os.Stat(fullPath)
		if os.IsNotExist(err) || stat.IsDir() {
			pLogger.PushWarn(logBuilder.WithMessage("file_not_found"))
			api.Response(pW, http.StatusNotAcceptable, "failed: file not found")
			return
		}

		size := stat.Size()
		chunks := uint64(math.Ceil(float64(size) / hlf_settings.CChunkSize))
		if uint64(chunk) >= chunks {
			pLogger.PushWarn(logBuilder.WithMessage("chunk_number"))
			api.Response(pW, http.StatusNotAcceptable, "failed: chunk number")
			return
		}

		file, err := os.Open(fullPath)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("open_file"))
			api.Response(pW, http.StatusNotAcceptable, "failed: open file")
			return
		}
		defer file.Close()

		buf := make([]byte, hlf_settings.CChunkSize)

		chunkOffset := int64(chunk) * hlf_settings.CChunkSize
		nS, err := file.Seek(chunkOffset, io.SeekStart)
		if err != nil || nS != chunkOffset {
			pLogger.PushWarn(logBuilder.WithMessage("seek_file"))
			api.Response(pW, http.StatusNotAcceptable, "failed: seek file")
			return
		}

		nR, err := file.Read(buf)
		if err != nil || (uint64(chunk) != chunks-1 && nR != hlf_settings.CChunkSize) {
			pLogger.PushWarn(logBuilder.WithMessage("chunk_number"))
			api.Response(pW, http.StatusNotAcceptable, "failed: chunk number")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		api.Response(pW, http.StatusOK, buf[:nR])
	}
}
