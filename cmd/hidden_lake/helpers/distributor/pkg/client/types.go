package client

import (
	"context"

	"github.com/number571/go-peer/cmd/hidden_lake/helpers/distributor/pkg/config"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/request"
	"github.com/number571/go-peer/cmd/hidden_lake/service/pkg/response"
)

type IClient interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)
	DistributeRequest(context.Context, request.IRequest) (response.IResponse, error)
}

type IRequester interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)
	DistributeRequest(context.Context, request.IRequest) (response.IResponse, error)
}
