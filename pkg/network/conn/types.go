package conn

import (
	"context"
	"net"
	"time"

	"github.com/number571/go-peer/pkg/types"

	net_message "github.com/number571/go-peer/pkg/network/message"
)

type IConn interface {
	types.ICloser

	GetSettings() ISettings
	GetSocket() net.Conn

	GetVSettings() IVSettings
	SetVSettings(IVSettings)

	WriteMessage(context.Context, net_message.IMessage) error
	ReadMessage(context.Context, chan<- struct{}) (net_message.IMessage, error)
}

type ISettings interface {
	GetLimitMessageSizeBytes() uint64
	GetLimitVoidSizeBytes() uint64
	GetWorkSizeBits() uint64

	GetDialTimeout() time.Duration
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
	GetWaitReadTimeout() time.Duration
}

type IVSettings interface {
	GetNetworkKey() string
}
