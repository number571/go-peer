package anonymity

import (
	"time"

	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity/adapters"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/types"
)

type (
	IHandlerF func(INode, asymmetric.IPubKey, []byte, []byte) []byte
)

type ISettings interface {
	GetServiceName() string
	GetNetworkMask() uint64
	GetRetryEnqueue() uint64
	GetFetchTimeWait() time.Duration
}

type INode interface {
	types.ICommand

	GetSettings() ISettings
	GetWrapperDB() IWrapperDB
	GetNetworkNode() network.INode
	GetMessageQueue() queue.IMessageQueue
	GetListPubKeys() asymmetric.IListPubKeys
	GetLogger() logger.ILogger

	HandleFunc(uint32, IHandlerF) INode
	HandleMessage(message.IMessage) // in runtime

	BroadcastPayload(asymmetric.IPubKey, adapters.IPayload) error
	FetchPayload(asymmetric.IPubKey, adapters.IPayload) ([]byte, error)
}

type IWrapperDB interface {
	types.ICloser

	Get() database.IKeyValueDB
	Set(database.IKeyValueDB) IWrapperDB
}
