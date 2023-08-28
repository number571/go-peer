package conn

import (
	"net"
	"time"

	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/types"
)

type IConn interface {
	types.ICloser

	GetSettings() ISettings
	GetSocket() net.Conn

	WritePayload(payload.IPayload) error
	ReadPayload() (payload.IPayload, error)
}

type ISettings interface {
	// for subsequent inheritance on multiple connections
	SetNetworkKey(string)

	GetNetworkKey() string
	GetMessageSizeBytes() uint64
	GetLimitVoidSize() uint64
	GetWaitReadDeadline() time.Duration
	GetReadDeadline() time.Duration
	GetWriteDeadline() time.Duration
}
