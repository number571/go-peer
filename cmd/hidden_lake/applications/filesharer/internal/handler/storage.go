package handler

import (
	"context"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/config"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/internal/stream"
	hlf_client "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/client"
	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/web"
	hls_client "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/client"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
)

type sStorage struct {
	*sTemplate
	FPage      uint64
	FAliasName string
	FFilesList []hlf_settings.SFileInfo
}

func StoragePage(
	pCtx context.Context,
	pLogger logger.ILogger,
	pCfg config.IConfig,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlf_settings.CServiceName, pR)

		if pR.URL.Path != "/friends/storage" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		query := pR.URL.Query()
		aliasName := query.Get("alias_name")
		if aliasName == "" {
			ErrorPage(pLogger, pCfg, "alias_name_error", "incorrect alias name")(pW, pR)
			return
		}

		hlsClient := getHLSClient(pCfg)
		if fileName := query.Get("file_name"); fileName != "" {
			downloadFile(pCtx, pLogger, pCfg, pW, pR, hlsClient)
			return
		}

		page, err := strconv.Atoi(query.Get("page"))
		if err != nil {
			page = 0
		}

		hlfClient := hlf_client.NewClient(
			hlf_client.NewBuilder(),
			hlf_client.NewRequester(hlsClient),
		)

		filesList, err := hlfClient.GetListFiles(pCtx, aliasName, uint64(page))
		if err != nil {
			ErrorPage(pLogger, pCfg, "get_files_list", "failed get list of files")(pW, pR)
			return
		}

		result := sStorage{
			sTemplate:  getTemplate(pCfg),
			FPage:      uint64(page),
			FAliasName: aliasName,
			FFilesList: filesList,
		}

		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"storage.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = t.Execute(pW, result)
	}
}

func downloadFile(
	pCtx context.Context,
	pLogger logger.ILogger,
	pCfg config.IConfig,
	pW http.ResponseWriter,
	pR *http.Request,
	pHlsClient hls_client.IClient,
) {
	query := pR.URL.Query()

	aliasName := query.Get("alias_name")
	fileName := query.Get("file_name")

	fileHash := query.Get("file_hash")
	if fileHash == "" || len(encoding.HexDecode(fileHash)) != hashing.CSHA256Size {
		ErrorPage(pLogger, pCfg, "file_hash_error", "incorrect file hash")(pW, pR)
		return
	}

	fileSize, err := strconv.ParseUint(query.Get("file_size"), 10, 64)
	if err != nil {
		ErrorPage(pLogger, pCfg, "file_size_error", "incorrect file size")(pW, pR)
		return
	}

	pW.Header().Set("Content-Type", "application/octet-stream")
	pW.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(fileName))
	pW.Header().Set("Content-Length", strconv.FormatUint(fileSize, 10))

	fileinfo := stream.NewFileInfo(fileName, fileHash, fileSize)
	stream, err := stream.BuildStream(pCtx, 2, pHlsClient, aliasName, fileinfo)
	if err != nil {
		ErrorPage(pLogger, pCfg, "build_stream", "build stream")(pW, pR)
		return
	}

	http.ServeContent(pW, pR, fileName, time.Now(), stream)
}
