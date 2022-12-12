package conn

import (
	"net"
	"time"

	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/go-peer/pkg/types"
)

type ISettings interface {
	GetNetworkKey() string
	GetMessageSize() uint64
	GetTimeWait() time.Duration
}

type IConn interface {
	Settings() ISettings
	types.ICloser

	Socket() net.Conn
	Request(pld payload.IPayload) (payload.IPayload, error)

	Write(payload.IPayload) error
	Read() payload.IPayload
}
