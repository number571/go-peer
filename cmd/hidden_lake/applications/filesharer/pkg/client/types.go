package client

import (
	"context"

	hlf_settings "github.com/number571/go-peer/cmd/hidden_lake/applications/filesharer/pkg/settings"
	hls_request "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
)

type IClient interface {
	GetListFiles(context.Context, string, uint64) ([]hlf_settings.SFileInfo, error)
	LoadFileChunk(context.Context, string, string, uint64) ([]byte, error)
}

type IRequester interface {
	GetListFiles(context.Context, string, hls_request.IRequest) ([]hlf_settings.SFileInfo, error)
	LoadFileChunk(context.Context, string, hls_request.IRequest) ([]byte, error)
}

type IBuilder interface {
	GetListFiles(uint64) hls_request.IRequest
	LoadFileChunk(string, uint64) hls_request.IRequest
}
