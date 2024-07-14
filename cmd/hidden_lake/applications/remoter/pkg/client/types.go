package client

import (
	"context"

	hls_request "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
)

type IClient interface {
	Exec(context.Context, string, ...string) ([]byte, error)
}

type IRequester interface {
	Exec(context.Context, string, hls_request.IRequest) ([]byte, error)
}

type IBuilder interface {
	Exec(...string) hls_request.IRequest
}
