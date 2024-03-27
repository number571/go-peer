package client

import (
	"context"

	hls_request "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
)

type IClient interface {
	PushMessage(context.Context, string, string, []byte) error
}

type IRequester interface {
	PushMessage(context.Context, string, hls_request.IRequest) error
}

type IBuilder interface {
	PushMessage(string, []byte) hls_request.IRequest
}
