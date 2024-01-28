package handler

import (
	"fmt"
	"html/template"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/applications/filer/internal/config"
	hlf_client "github.com/number571/go-peer/cmd/hidden_lake/applications/filer/pkg/client"
	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filer/pkg/settings"
	"github.com/number571/go-peer/cmd/hidden_lake/applications/filer/web"
	http_logger "github.com/number571/go-peer/internal/logger/http"
	"github.com/number571/go-peer/pkg/logger"
)

type sStorage struct {
	*sTemplate
	FName      string
	FPage      uint64
	FAliasName string
	FFilesList []hlf_settings.SFileInfo
}

func StoragePage(pLogger logger.ILogger, pCfg config.IConfig, pPathTo string) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlf_settings.CServiceName, pR)

		if pR.URL.Path != "/friends/storage" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		hlsClient := getClient(pCfg)
		hlfClient := hlf_client.NewClient(
			hlf_client.NewBuilder(),
			hlf_client.NewRequester(hlsClient),
		)

		query := pR.URL.Query()
		aliasName := query.Get("alias_name")
		if aliasName == "" {
			ErrorPage(pLogger, pCfg, "alias_name_error", "incorrect alias name")(pW, pR)
			return
		}

		fileName := ""
		pR.ParseForm()

		switch pR.FormValue("method") {
		case http.MethodPost:
			fileName = pR.FormValue("file_name")
			if fileName == "" {
				ErrorPage(pLogger, pCfg, "file_name_error", "incorrect file name")(pW, pR)
				return
			}

			fileHash := pR.FormValue("file_hash")
			if fileHash == "" {
				ErrorPage(pLogger, pCfg, "file_hash_error", "incorrect file hash")(pW, pR)
				return
			}

			fileSize, err := strconv.Atoi(pR.FormValue("file_size"))
			if err != nil {
				ErrorPage(pLogger, pCfg, "file_size_error", "incorrect file size")(pW, pR)
				return
			}

			baseFileName := getUniqFileName(pPathTo, fileName)
			tempFileName := baseFileName + ".tmp"

			file, err := os.Create(tempFileName)
			if err != nil {
				ErrorPage(pLogger, pCfg, "create_file_error", "failed to create file")(pW, pR)
				return
			}
			defer func() {
				file.Close()
				os.Remove(tempFileName)
			}()

			gotSize := 0
			retryNum := 3

			chunksCount := uint64(math.Ceil(float64(fileSize) / hlf_settings.CChunkSize))
			for i := uint64(0); i < chunksCount; i++ {
				for j := 1; j <= retryNum; j++ {
					chunk, err := hlfClient.LoadFileChunk(aliasName, fileName, i)
					if err != nil {
						if j == retryNum {
							ErrorPage(pLogger, pCfg, "load_chunk_error", "failed to load chunk")(pW, pR)
							return
						}
						continue
					}
					n, err := file.Write(chunk)
					if err != nil {
						ErrorPage(pLogger, pCfg, "write_to_file", "failed write to file")(pW, pR)
						return
					}
					if i != chunksCount-1 && n != hlf_settings.CChunkSize {
						ErrorPage(pLogger, pCfg, "invalid_chunk_size", "got invalid chunk size")(pW, pR)
						return
					}
					gotSize += n
					break
				}
			}

			if gotSize != fileSize {
				ErrorPage(pLogger, pCfg, "size_file", "invalid size file")(pW, pR)
				return
			}

			if getFileHash(tempFileName) != fileHash {
				ErrorPage(pLogger, pCfg, "hash_file", "invalid hash file")(pW, pR)
				return
			}

			// baseFile <- tempFile
			if err := copyFile(baseFileName, tempFileName); err != nil {
				ErrorPage(pLogger, pCfg, "copy_file", "failed copy file")(pW, pR)
				return
			}
		}

		page, err := strconv.Atoi(query.Get("page"))
		if err != nil {
			page = 0
		}

		filesList, err := hlfClient.GetListFiles(aliasName, uint64(page))
		if err != nil {
			ErrorPage(pLogger, pCfg, "get_files_list", "failed get list of files")(pW, pR)
			return
		}

		result := sStorage{
			sTemplate:  getTemplate(pCfg),
			FPage:      uint64(page),
			FName:      fileName,
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
		t.Execute(pW, result)
	}
}

func getUniqFileName(pPathTo, pFileName string) string {
	origName := fmt.Sprintf("%s/%s/%s", pPathTo, hlf_settings.CPathLoadedSTG, pFileName)
	if _, err := os.Stat(origName); os.IsNotExist(err) {
		return origName
	}
	return fmt.Sprintf(
		"%s/%s/%s_%s",
		pPathTo,
		hlf_settings.CPathLoadedSTG,
		time.Now().Format("20060102150405"),
		pFileName,
	)
}

func copyFile(dst, src string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
