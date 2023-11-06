package conn

import (
	"net"
	"time"

	"github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/types"
)

type IConn interface {
	types.ICloser

	GetSettings() ISettings
	GetSocket() net.Conn

	WriteMessage(message.IMessage) error
	ReadMessage(chan struct{}) (message.IMessage, error)
}

type ISettings interface {
	message.ISettings

	// for subsequent inheritance on multiple connections
	SetNetworkKey(string)

	GetMessageSizeBytes() uint64
	GetLimitVoidSize() uint64
	GetWaitReadDeadline() time.Duration
	GetReadDeadline() time.Duration
	GetWriteDeadline() time.Duration
}
