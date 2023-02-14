package testutils

import (
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/client/queue"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/friends"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/database"
	testutils "github.com/number571/go-peer/test/_data"
)

const (
	TCMessageSize = uint64(100 << 10)
	TCWorkSize    = 10
)

func TestNewNode(pathDB string) anonymity.INode {
	node := anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FTimeWait: 30 * time.Second,
		}),
		logger.NewLogger(logger.NewSettings(&logger.SSettings{})),
		database.NewLevelDB(
			database.NewSettings(&database.SSettings{
				FPath:      pathDB,
				FHashing:   true,
				FCipherKey: []byte(testutils.TcKey1),
			}),
		),
		TestNewNetworkNode(),
		queue.NewQueue(
			queue.NewSettings(&queue.SSettings{
				FCapacity:     10,
				FPullCapacity: 5,
				FDuration:     500 * time.Millisecond,
			}),
			client.NewClient(
				client.NewSettings(&client.SSettings{
					FWorkSize:    TCWorkSize,
					FMessageSize: TCMessageSize,
				}),
				asymmetric.LoadRSAPrivKey(testutils.TcPrivKey),
			),
		),
		friends.NewF2F(),
	)
	return node
}

func TestNewNetworkNode() network.INode {
	return network.NewNode(
		network.NewSettings(&network.SSettings{
			FCapacity:    (1 << 10),
			FMaxConnects: 10,
			FConnSettings: conn.NewSettings(&conn.SSettings{
				FMessageSize: TCMessageSize,
				FTimeWait:    5 * time.Second,
			}),
		}),
	)
}
