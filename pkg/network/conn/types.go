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

	WriteMessage(context.Context, net_message.IMessage) error
	ReadMessage(context.Context, chan<- struct{}) (net_message.IMessage, error)
}

type ISettings interface {
	net_message.ISettings
	GetLimitMessageSizeBytes() uint64
	GetDialTimeout() time.Duration
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
	GetWaitReadTimeout() time.Duration
}
