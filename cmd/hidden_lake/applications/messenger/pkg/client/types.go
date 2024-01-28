package client

import (
	hls_request "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
)

type IClient interface {
	PushMessage(string, string, string, []byte) error
}

type IRequester interface {
	PushMessage(string, hls_request.IRequest) error
}

type IBuilder interface {
	PushMessage(string, string, []byte) hls_request.IRequest
}
