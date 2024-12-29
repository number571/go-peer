package conn

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/number571/go-peer/pkg/message/layer1"
)

type IConn interface {
	io.Closer

	GetSettings() ISettings
	GetSocket() net.Conn

	WriteMessage(context.Context, layer1.IMessage) error
	ReadMessage(context.Context, chan<- struct{}) (layer1.IMessage, error)
}

type ISettings interface {
	GetMessageSettings() layer1.ISettings
	GetLimitMessageSizeBytes() uint64
	GetDialTimeout() time.Duration
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
	GetWaitReadTimeout() time.Duration
}
