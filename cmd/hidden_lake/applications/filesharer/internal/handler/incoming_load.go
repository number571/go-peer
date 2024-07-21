package handler

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/utils"
	"github.com/number571/go-peer/internal/api"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"

	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/settings"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
)

func HandleIncomingLoadHTTP(
	pCtx context.Context,
	pLogger logger.ILogger,
	pStgPath string,
	pHlsClient hls_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hls_settings.CHeaderResponseMode, hls_settings.CHeaderResponseModeON)

		logBuilder := http_logger.NewLogBuilder(hlf_settings.CServiceName, pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		query := pR.URL.Query()

		name := filepath.Base(query.Get("name"))
		if name != query.Get("name") {
			pLogger.PushWarn(logBuilder.WithMessage("got_another_name"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: got another name")
			return
		}

		chunk, err := strconv.Atoi(query.Get("chunk"))
		if err != nil || chunk < 0 {
			pLogger.PushWarn(logBuilder.WithMessage("incorrect_chunk"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: incorrect chunk")
			return
		}

		fullPath := filepath.Join(pStgPath, name)
		stat, err := os.Stat(fullPath)
		if os.IsNotExist(err) || stat.IsDir() {
			pLogger.PushWarn(logBuilder.WithMessage("file_not_found"))
			_ = api.Response(pW, http.StatusNotFound, "failed: file not found")
			return
		}

		chunkSize, err := utils.GetMessageLimit(pCtx, pHlsClient)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_chunk_size"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: get chunk size")
			return
		}

		chunks := utils.GetChunksCount(uint64(stat.Size()), chunkSize)
		if uint64(chunk) >= chunks {
			pLogger.PushWarn(logBuilder.WithMessage("chunk_number"))
			_ = api.Response(pW, http.StatusLengthRequired, "failed: chunk number")
			return
		}

		file, err := os.Open(fullPath)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("open_file"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: open file")
			return
		}
		defer file.Close()

		buf := make([]byte, chunkSize)
		chunkOffset := int64(chunk) * int64(chunkSize)

		nS, err := file.Seek(chunkOffset, io.SeekStart)
		if err != nil || nS != chunkOffset {
			pLogger.PushWarn(logBuilder.WithMessage("seek_file"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: seek file")
			return
		}

		nR, err := file.Read(buf)
		if err != nil || (uint64(chunk) != chunks-1 && uint64(nR) != chunkSize) {
			pLogger.PushWarn(logBuilder.WithMessage("chunk_number"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: chunk number")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, buf[:nR])
	}
}
