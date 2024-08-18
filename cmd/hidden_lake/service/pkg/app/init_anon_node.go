package app

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/number571/go-peer/cmd/hidden_lake/service/internal/handler"
	"github.com/number571/go-peer/pkg/client/message"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/network"
	"github.com/number571/go-peer/pkg/network/anonymity"
	"github.com/number571/go-peer/pkg/network/anonymity/queue"
	"github.com/number571/go-peer/pkg/network/conn"
	"github.com/number571/go-peer/pkg/storage/cache/lru"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/go-peer/pkg/utils"

	hls_settings "github.com/number571/go-peer/cmd/hidden_lake/service/pkg/settings"
	"github.com/number571/go-peer/pkg/client"
)

const (
	cPPKKey = "__PSD_PUB_KEY__"
)

func init() {
	// The hls.db file stores only message hashes.
	// If the other parameters have the same size as the hashes,
	// then there will be a chance of overwriting.
	if len(cPPKKey) == hashing.CSHA256Size {
		panic("len(cPPKKey) == hashing.CSHA256Size")
	}
}

func (p *sApp) initAnonNode() error {
	var (
		cfg         = p.fCfgW.GetConfig()
		cfgSettings = cfg.GetSettings()
	)

	kvDatabase, err := database.NewKVDatabase(
		database.NewSettings(&database.SSettings{
			FPath: filepath.Join(p.fPathTo, hls_settings.CPathDB),
		}),
	)
	if err != nil {
		return utils.MergeErrors(ErrOpenKVDatabase, err)
	}

	psdPubKey, err := getPsdPubKey(kvDatabase, cfgSettings.GetKeySizeBits())
	if err != nil {
		return utils.MergeErrors(ErrGetPsdPubKey, err)
	}

	client := client.NewClient(
		message.NewSettings(&message.SSettings{
			FMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
			FKeySizeBits:      p.fPrivKey.GetSize(),
		}),
		p.fPrivKey,
	)
	if client.GetMessageLimit() <= encoding.CSizeUint64 {
		return utils.MergeErrors(ErrMessageSizeLimit, err)
	}

	p.fNode = anonymity.NewNode(
		anonymity.NewSettings(&anonymity.SSettings{
			FServiceName:  hls_settings.CServiceName,
			FF2FDisabled:  cfgSettings.GetF2FDisabled(),
			FNetworkMask:  hls_settings.CNetworkMask,
			FFetchTimeout: time.Duration(cfgSettings.GetFetchTimeoutMS()) * time.Millisecond,
		}),
		// Insecure to use logging in real anonymity projects!
		// Logging should only be used in overview or testing;
		p.fAnonLogger,
		kvDatabase,
		network.NewNode(
			network.NewSettings(&network.SSettings{
				FAddress:      cfg.GetAddress().GetTCP(),
				FMaxConnects:  hls_settings.CNetworkMaxConns,
				FReadTimeout:  hls_settings.CNetworkReadTimeout,
				FWriteTimeout: hls_settings.CNetworkWriteTimeout,
				FConnSettings: conn.NewSettings(&conn.SSettings{
					FLimitMessageSizeBytes: cfgSettings.GetMessageSizeBytes() + cfgSettings.GetRandMessageSizeBytes(),
					FWorkSizeBits:          cfgSettings.GetWorkSizeBits(),
					FWaitReadTimeout:       hls_settings.CConnWaitReadTimeout,
					FDialTimeout:           hls_settings.CConnDialTimeout,
					FReadTimeout:           hls_settings.CNetworkReadTimeout,
					FWriteTimeout:          hls_settings.CNetworkWriteTimeout,
				}),
			}),
			conn.NewVSettings(&conn.SVSettings{
				FNetworkKey: cfgSettings.GetNetworkKey(),
			}),
			lru.NewLRUCache(hls_settings.CNetworkQueueCapacity),
		),
		queue.NewQBTaskProcessor(
			queue.NewSettings(&queue.SSettings{
				FNetworkMask:          hls_settings.CNetworkMask,
				FWorkSizeBits:         cfgSettings.GetWorkSizeBits(),
				FMainPoolCapacity:     hls_settings.CQueueMainPoolCapacity,
				FRandPoolCapacity:     hls_settings.CQueueRandPoolCapacity,
				FParallel:             p.fParallel,
				FRandMessageSizeBytes: cfgSettings.GetRandMessageSizeBytes(),
				FQueuePeriod:          time.Duration(cfgSettings.GetQueuePeriodMS()) * time.Millisecond,
				FRandQueuePeriod:      time.Duration(cfgSettings.GetRandQueuePeriodMS()) * time.Millisecond,
			}),
			queue.NewVSettings(&queue.SVSettings{
				FNetworkKey: cfgSettings.GetNetworkKey(),
			}),
			client,
			psdPubKey,
		),
		func() asymmetric.IListPubKeys {
			f2f := asymmetric.NewListPubKeys()
			for _, pubKey := range cfg.GetFriends() {
				f2f.AddPubKey(pubKey)
			}
			return f2f
		}(),
	).HandleFunc(
		hls_settings.CServiceMask,
		handler.HandleServiceTCP(p.fCfgW),
	)

	return nil
}

func getPsdPubKey(pDB database.IKVDatabase, pKeySize uint64) (asymmetric.IPubKey, error) {
	ppk, err := pDB.Get([]byte(cPPKKey))
	if err == nil {
		pubKey := asymmetric.LoadRSAPubKey(ppk)
		if pubKey == nil {
			return nil, ErrInvalidPsdPubKey
		}
		return pubKey, nil
	}
	if errors.Is(err, database.ErrNotFound) {
		pubKey := asymmetric.NewRSAPrivKey(pKeySize).GetPubKey()
		if err := pDB.Set([]byte(cPPKKey), pubKey.ToBytes()); err != nil {
			return nil, ErrSetPsdPubKey
		}
		return pubKey, nil
	}
	return nil, utils.MergeErrors(err, ErrReadKVDatabase)
}
